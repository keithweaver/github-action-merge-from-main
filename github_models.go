package main

import "time"

// PullRequest represents a GitHub pull request
type PullRequest struct {
	URL                string     `json:"url"`
	ID                 int        `json:"id"`
	NodeID             string     `json:"node_id"`
	HTMLURL            string     `json:"html_url"`
	DiffURL            string     `json:"diff_url"`
	PatchURL           string     `json:"patch_url"`
	IssueURL           string     `json:"issue_url"`
	CommitsURL         string     `json:"commits_url"`
	ReviewCommentsURL  string     `json:"review_comments_url"`
	ReviewCommentURL   string     `json:"review_comment_url"`
	CommentsURL        string     `json:"comments_url"`
	StatusesURL        string     `json:"statuses_url"`
	Number             int        `json:"number"`
	State              string     `json:"state"`
	Locked             bool       `json:"locked"`
	Title              string     `json:"title"`
	User               User       `json:"user"`
	Body               string     `json:"body"`
	Labels             []Label    `json:"labels"`
	Milestone          *Milestone `json:"milestone"`
	ActiveLockReason   string     `json:"active_lock_reason"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	ClosedAt           *time.Time `json:"closed_at"`
	MergedAt           *time.Time `json:"merged_at"`
	MergeCommitSha     string     `json:"merge_commit_sha"`
	Assignee           *User      `json:"assignee"`
	Assignees          []User     `json:"assignees"`
	RequestedReviewers []User     `json:"requested_reviewers"`
	RequestedTeams     []Team     `json:"requested_teams"`
	Head               Branch     `json:"head"`
	Base               Branch     `json:"base"`
	Links              Links      `json:"_links"`
	AuthorAssociation  string     `json:"author_association"`
	AutoMerge          *AutoMerge `json:"auto_merge"`
	Draft              bool       `json:"draft"`
}

// User represents a GitHub user
type User struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}

// Label represents a GitHub label
type Label struct {
	ID          int    `json:"id"`
	NodeID      string `json:"node_id"`
	URL         string `json:"url"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

// Milestone represents a GitHub milestone
type Milestone struct {
	URL          string     `json:"url"`
	HTMLURL      string     `json:"html_url"`
	LabelsURL    string     `json:"labels_url"`
	ID           int        `json:"id"`
	NodeID       string     `json:"node_id"`
	Number       int        `json:"number"`
	State        string     `json:"state"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Creator      User       `json:"creator"`
	OpenIssues   int        `json:"open_issues"`
	ClosedIssues int        `json:"closed_issues"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	ClosedAt     *time.Time `json:"closed_at"`
	DueOn        *time.Time `json:"due_on"`
}

// Team represents a GitHub team
type Team struct {
	ID                  int    `json:"id"`
	NodeID              string `json:"node_id"`
	URL                 string `json:"url"`
	HTMLURL             string `json:"html_url"`
	Name                string `json:"name"`
	Slug                string `json:"slug"`
	Description         string `json:"description"`
	Privacy             string `json:"privacy"`
	Permission          string `json:"permission"`
	NotificationSetting string `json:"notification_setting"`
	MembersURL          string `json:"members_url"`
	RepositoriesURL     string `json:"repositories_url"`
	Parent              *Team  `json:"parent"`
}

// Branch represents a GitHub branch reference
type Branch struct {
	Label string     `json:"label"`
	Ref   string     `json:"ref"`
	Sha   string     `json:"sha"`
	User  User       `json:"user"`
	Repo  Repository `json:"repo"`
}

// Repository represents a GitHub repository
type Repository struct {
	ID                  int         `json:"id"`
	NodeID              string      `json:"node_id"`
	Name                string      `json:"name"`
	FullName            string      `json:"full_name"`
	Owner               User        `json:"owner"`
	Private             bool        `json:"private"`
	HTMLURL             string      `json:"html_url"`
	Description         string      `json:"description"`
	Fork                bool        `json:"fork"`
	URL                 string      `json:"url"`
	ArchiveURL          string      `json:"archive_url"`
	AssigneesURL        string      `json:"assignees_url"`
	BlobsURL            string      `json:"blobs_url"`
	BranchesURL         string      `json:"branches_url"`
	CollaboratorsURL    string      `json:"collaborators_url"`
	CommentsURL         string      `json:"comments_url"`
	CommitsURL          string      `json:"commits_url"`
	CompareURL          string      `json:"compare_url"`
	ContentsURL         string      `json:"contents_url"`
	ContributorsURL     string      `json:"contributors_url"`
	DeploymentsURL      string      `json:"deployments_url"`
	DownloadsURL        string      `json:"downloads_url"`
	EventsURL           string      `json:"events_url"`
	ForksURL            string      `json:"forks_url"`
	GitCommitsURL       string      `json:"git_commits_url"`
	GitRefsURL          string      `json:"git_refs_url"`
	GitTagsURL          string      `json:"git_tags_url"`
	GitURL              string      `json:"git_url"`
	IssueCommentURL     string      `json:"issue_comment_url"`
	IssueEventsURL      string      `json:"issue_events_url"`
	IssuesURL           string      `json:"issues_url"`
	KeysURL             string      `json:"keys_url"`
	LabelsURL           string      `json:"labels_url"`
	LanguagesURL        string      `json:"languages_url"`
	MergesURL           string      `json:"merges_url"`
	MilestonesURL       string      `json:"milestones_url"`
	NotificationsURL    string      `json:"notifications_url"`
	PullsURL            string      `json:"pulls_url"`
	ReleasesURL         string      `json:"releases_url"`
	SSHURL              string      `json:"ssh_url"`
	StargazersURL       string      `json:"stargazers_url"`
	StatusesURL         string      `json:"statuses_url"`
	SubscribersURL      string      `json:"subscribers_url"`
	SubscriptionURL     string      `json:"subscription_url"`
	TagsURL             string      `json:"tags_url"`
	TeamsURL            string      `json:"teams_url"`
	TreesURL            string      `json:"trees_url"`
	CloneURL            string      `json:"clone_url"`
	MirrorURL           string      `json:"mirror_url"`
	HooksURL            string      `json:"hooks_url"`
	SvnURL              string      `json:"svn_url"`
	Homepage            string      `json:"homepage"`
	Language            *string     `json:"language"`
	ForksCount          int         `json:"forks_count"`
	StargazersCount     int         `json:"stargazers_count"`
	WatchersCount       int         `json:"watchers_count"`
	Size                int         `json:"size"`
	DefaultBranch       string      `json:"default_branch"`
	OpenIssuesCount     int         `json:"open_issues_count"`
	IsTemplate          bool        `json:"is_template"`
	Topics              []string    `json:"topics"`
	HasIssues           bool        `json:"has_issues"`
	HasProjects         bool        `json:"has_projects"`
	HasWiki             bool        `json:"has_wiki"`
	HasPages            bool        `json:"has_pages"`
	HasDownloads        bool        `json:"has_downloads"`
	Archived            bool        `json:"archived"`
	Disabled            bool        `json:"disabled"`
	Visibility          string      `json:"visibility"`
	PushedAt            time.Time   `json:"pushed_at"`
	CreatedAt           time.Time   `json:"created_at"`
	UpdatedAt           time.Time   `json:"updated_at"`
	Permissions         Permission  `json:"permissions"`
	AllowRebaseMerge    bool        `json:"allow_rebase_merge"`
	TemplateRepository  *Repository `json:"template_repository"`
	TempCloneToken      string      `json:"temp_clone_token"`
	AllowSquashMerge    bool        `json:"allow_squash_merge"`
	AllowAutoMerge      bool        `json:"allow_auto_merge"`
	DeleteBranchOnMerge bool        `json:"delete_branch_on_merge"`
	AllowMergeCommit    bool        `json:"allow_merge_commit"`
	SubscribersCount    int         `json:"subscribers_count"`
	NetworkCount        int         `json:"network_count"`
	License             *License    `json:"license"`
	Forks               int         `json:"forks"`
	OpenIssues          int         `json:"open_issues"`
	Watchers            int         `json:"watchers"`
}

// Permission represents repository permissions
type Permission struct {
	Admin bool `json:"admin"`
	Push  bool `json:"push"`
	Pull  bool `json:"pull"`
}

// License represents a repository license
type License struct {
	Key     string `json:"key"`
	Name    string `json:"name"`
	URL     string `json:"url"`
	SpdxID  string `json:"spdx_id"`
	NodeID  string `json:"node_id"`
	HTMLURL string `json:"html_url"`
}

// Links represents the _links object in a PR
type Links struct {
	Self           Link `json:"self"`
	HTML           Link `json:"html"`
	Issue          Link `json:"issue"`
	Comments       Link `json:"comments"`
	ReviewComments Link `json:"review_comments"`
	ReviewComment  Link `json:"review_comment"`
	Commits        Link `json:"commits"`
	Statuses       Link `json:"statuses"`
}

// Link represents a single link
type Link struct {
	Href string `json:"href"`
}

// AutoMerge represents the auto_merge setting
type AutoMerge struct {
	EnabledBy     User   `json:"enabled_by"`
	MergeMethod   string `json:"merge_method"`
	CommitTitle   string `json:"commit_title"`
	CommitMessage string `json:"commit_message"`
}

// MergeRequest represents the request body for merging a PR
type MergeRequest struct {
	CommitTitle   string `json:"commit_title,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
	Sha           string `json:"sha,omitempty"`
	MergeMethod   string `json:"merge_method,omitempty"`
}

// MergeResponse represents the response from merging a PR
type MergeResponse struct {
	Sha     string `json:"sha"`
	Merged  bool   `json:"merged"`
	Message string `json:"message"`
}

type createPullRequest struct {
	Title string `json:"title"`
	Head  string `json:"head"`
	Base  string `json:"base"`
	Body  string `json:"body"`
}

// CombinedStatus represents the combined status for a commit.
type CombinedStatus struct {
	State      string         `json:"state"`
	Statuses   []Status       `json:"statuses"`
	SHA        string         `json:"sha"`
	TotalCount int            `json:"total_count"`
	CommitURL  string         `json:"commit_url"`
	URL        string         `json:"url"`
	Repository RepositoryMeta `json:"repository"`
}

// Status represents individual status checks.
type Status struct {
	State       string    `json:"state"`
	TargetURL   string    `json:"target_url"`
	Description string    `json:"description"`
	Context     string    `json:"context"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// RepositoryMeta is a minimal repository representation for status responses.
type RepositoryMeta struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Owner User   `json:"owner"`
}
