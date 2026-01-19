package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type config struct {
	AccessToken    string
	CommitPrefix   string
	Commands       []string
	RepoOwner      string
	RepoName       string
	BaseBranch     string
	IgnorePrefixes []string
	RunOnPrefixes  []string
	RunOnContains  []string
	WaitSeconds    int
	CIWaitTimeout  time.Duration
	CIWaitInterval time.Duration
	PushRemote     string
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		fail(err)
	}

	shouldRun, err := ConfirmShouldRun(cfg)
	if err != nil {
		fail(err)
	}

	if !shouldRun {
		log.Println("Last commit uses an ignore prefix. Exiting without action.")
		return
	}

	if err := RunCommands(cfg.Commands); err != nil {
		fail(err)
	}

	changed, err := hasChanges()
	if err != nil {
		fail(err)
	}

	if !changed {
		log.Println("No changes detected after running commands. Nothing to commit.")
		return
	}

	client := NewGitHubClient(cfg.AccessToken, cfg.RepoOwner, cfg.RepoName)

	pr, headSHA, err := CommitAndOpenPR(cfg, client)
	if err != nil {
		fail(err)
	}

	Wait(cfg)

	if err := WaitForCI(cfg, client, headSHA); err != nil {
		fail(err)
	}

	if err := Merge(cfg, client, pr); err != nil {
		fail(err)
	}

	log.Println("Completed merge from main.")
}

// ConfirmShouldRun returns false when the last commit is already an auto-merge.
func ConfirmShouldRun(cfg config) (bool, error) {
	msg, err := latestCommitMessage()
	if err != nil {
		return false, err
	}

	for _, prefix := range cfg.IgnorePrefixes {
		if prefix == "" {
			continue
		}
		if strings.HasPrefix(msg, prefix) {
			return false, nil
		}
	}

	// Run-on filters: if provided, they must match.
	prefixMatch := len(cfg.RunOnPrefixes) == 0
	for _, prefix := range cfg.RunOnPrefixes {
		if prefix == "" {
			continue
		}
		if strings.HasPrefix(msg, prefix) {
			prefixMatch = true
			break
		}
	}

	containsMatch := len(cfg.RunOnContains) == 0
	for _, marker := range cfg.RunOnContains {
		if marker == "" {
			continue
		}
		if strings.Contains(msg, marker) {
			containsMatch = true
			break
		}
	}

	return prefixMatch && containsMatch, nil
}

// RunCommands executes the provided commands sequentially.
func RunCommands(commands []string) error {
	for _, cmd := range commands {
		command := strings.TrimSpace(cmd)
		if command == "" {
			continue
		}
		log.Printf("Running command: %s\n", command)
		c := exec.Command("bash", "-lc", command)
		c.Stdout = os.Stdout
		c.Stderr = os.Stderr
		if err := c.Run(); err != nil {
			return fmt.Errorf("command failed (%s): %w", command, err)
		}
	}
	return nil
}

// CommitAndOpenPR commits changes, pushes a branch, and opens a PR.
func CommitAndOpenPR(cfg config, client *GitHubClient) (*PullRequest, string, error) {
	branchName := fmt.Sprintf("auto-merge-%d", time.Now().Unix())
	if err := runGit("checkout", "-b", branchName); err != nil {
		return nil, "", fmt.Errorf("failed to create branch: %w", err)
	}

	if err := ensureGitUser(); err != nil {
		return nil, "", fmt.Errorf("failed to configure git user: %w", err)
	}

	if err := runGit("add", "--all"); err != nil {
		return nil, "", fmt.Errorf("failed to add changes: %w", err)
	}

	commitMessage := fmt.Sprintf("%s Merge from %s", cfg.CommitPrefix, cfg.BaseBranch)
	if err := runGit("commit", "-m", commitMessage); err != nil {
		return nil, "", fmt.Errorf("failed to commit changes: %w", err)
	}

	pushURL := cfg.PushRemote
	if pushURL == "" {
		pushURL = fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s.git", cfg.AccessToken, cfg.RepoOwner, cfg.RepoName)
	}
	if err := runGit("push", pushURL, branchName); err != nil {
		return nil, "", fmt.Errorf("failed to push branch: %w", err)
	}

	title := commitMessage
	body := "Automated updates from main."
	pr, err := client.CreatePullRequest(title, branchName, cfg.BaseBranch, body)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create pull request: %w", err)
	}

	sha, err := gitHeadSHA()
	if err != nil {
		return nil, "", err
	}

	return pr, sha, nil
}

// Wait pauses between PR creation and CI checks.
func Wait(cfg config) {
	wait := cfg.WaitSeconds
	if wait <= 0 {
		wait = 30
	}
	log.Printf("Waiting %d seconds before checking CI status...\n", wait)
	time.Sleep(time.Duration(wait) * time.Second)
}

// WaitForCI polls GitHub for combined status until success/failure or timeout.
func WaitForCI(cfg config, client *GitHubClient, sha string) error {
	timeout := cfg.CIWaitTimeout
	if timeout <= 0 {
		timeout = 15 * time.Minute
	}
	interval := cfg.CIWaitInterval
	if interval <= 0 {
		interval = 10 * time.Second
	}

	start := time.Now()
	for {
		status, err := client.GetCombinedStatus(sha)
		if err != nil {
			return fmt.Errorf("failed to fetch combined status: %w", err)
		}

		state := strings.ToLower(status.State)
		if state == "success" {
			return nil
		}
		if state == "failure" || state == "error" {
			return fmt.Errorf("ci reported %s", state)
		}
		if time.Since(start) > timeout {
			return fmt.Errorf("ci did not finish within %s", timeout)
		}

		log.Printf("CI status is %s; checking again in %s...\n", status.State, interval)
		time.Sleep(interval)
	}
}

// Merge completes the PR using squash merge.
func Merge(cfg config, client *GitHubClient, pr *PullRequest) error {
	commitTitle := pr.Title
	commitMessage := fmt.Sprintf("%s Squash merge by automation", cfg.CommitPrefix)

	merged, err := client.MergePullRequest(pr.Number, commitTitle, commitMessage)
	if err != nil {
		return err
	}
	if !merged {
		return errors.New("merge API returned false")
	}
	return nil
}

func loadConfig() (config, error) {
	token := strings.TrimSpace(firstNonEmpty(os.Getenv("INPUT_GITHUB_ACCESS_TOKEN"), os.Getenv("GITHUB_ACCESS_TOKEN")))
	if token == "" {
		return config{}, errors.New("github access token is required")
	}

	commitPrefix := strings.TrimSpace(os.Getenv("INPUT_COMMIT_PREFIX"))
	if commitPrefix == "" {
		commitPrefix = "[Auto Merge]"
	}

	commandsRaw := os.Getenv("INPUT_COMMANDS")
	commands := splitCommands(commandsRaw)
	if len(commands) == 0 {
		return config{}, errors.New("at least one command is required")
	}

	repo := os.Getenv("GITHUB_REPOSITORY")
	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return config{}, fmt.Errorf("invalid GITHUB_REPOSITORY: %s", repo)
	}

	baseBranch := os.Getenv("GITHUB_REF_NAME")
	if baseBranch == "" {
		baseBranch = "main"
	}

	prefixes := []string{"Auto Merge", "[Auto Merge]:", commitPrefix}
	extraPrefixes := parsePrefixes(os.Getenv("PREFIXES_TO_IGNORE"))
	prefixes = append(prefixes, extraPrefixes...)

	runPrefixes := parsePrefixes(os.Getenv("PREFIXES_TO_RUN_ON"))
	runContains := parsePrefixes(os.Getenv("CONTAINS_TO_RUN_ON"))

	return config{
		AccessToken:    token,
		CommitPrefix:   commitPrefix,
		Commands:       commands,
		RepoOwner:      parts[0],
		RepoName:       parts[1],
		BaseBranch:     baseBranch,
		IgnorePrefixes: prefixes,
		RunOnPrefixes:  runPrefixes,
		RunOnContains:  runContains,
		WaitSeconds:    30,
		CIWaitTimeout:  15 * time.Minute,
		CIWaitInterval: 10 * time.Second,
		PushRemote:     "",
	}, nil
}

func latestCommitMessage() (string, error) {
	out, err := exec.Command("git", "log", "-1", "--pretty=%s").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get latest commit message: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func hasChanges() (bool, error) {
	out, err := exec.Command("git", "status", "--porcelain").Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}
	return strings.TrimSpace(string(out)) != "", nil
}

func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func gitHeadSHA() (string, error) {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("failed to get head sha: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

func ensureGitUser() error {
	user := os.Getenv("GITHUB_ACTOR")
	if user == "" {
		user = "github-actions"
	}
	email := fmt.Sprintf("%s@users.noreply.github.com", user)
	if err := runGit("config", "user.name", user); err != nil {
		return err
	}
	return runGit("config", "user.email", email)
}

func splitCommands(raw string) []string {
	lines := strings.Split(raw, "\n")
	var cmds []string
	for _, line := range lines {
		for _, segment := range strings.Split(line, ",") {
			if c := strings.TrimSpace(segment); c != "" {
				cmds = append(cmds, c)
			}
		}
	}
	return cmds
}

func parsePrefixes(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	var prefixes []string
	for _, p := range parts {
		if s := strings.TrimSpace(p); s != "" {
			prefixes = append(prefixes, s)
		}
	}
	return prefixes
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func fail(err error) {
	log.Println(err)
	os.Exit(1)
}
