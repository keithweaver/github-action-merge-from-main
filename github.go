package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	githubAPIBaseURL = "https://api.github.com"
	perPage          = 100
)

// GitHubClient handles GitHub API requests
type GitHubClient struct {
	token     string
	repoOwner string
	repo      string
	client    *http.Client
}

// NewGitHubClient creates a new GitHub API client
func NewGitHubClient(token, repoOwner, repo string) *GitHubClient {
	return &GitHubClient{
		token:     token,
		repoOwner: repoOwner,
		repo:      repo,
		client:    &http.Client{},
	}
}

// CreatePullRequest opens a pull request from head to base.
func (c *GitHubClient) CreatePullRequest(title, head, base, body string) (*PullRequest, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls", githubAPIBaseURL, c.repoOwner, c.repo)
	reqBody := createPullRequest{
		Title: title,
		Head:  head,
		Base:  base,
		Body:  body,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	c.decorateHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &pr, nil
}

// GetCombinedStatus returns the combined status for a commit.
func (c *GitHubClient) GetCombinedStatus(sha string) (*CombinedStatus, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/commits/%s/status", githubAPIBaseURL, c.repoOwner, c.repo, sha)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	c.decorateHeaders(req)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var status CombinedStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}
	return &status, nil
}

// MergePullRequest merges a pull request using squash strategy.
func (c *GitHubClient) MergePullRequest(prNumber int, prTitle, prMessage string) (bool, error) {
	url := fmt.Sprintf("%s/repos/%s/%s/pulls/%d/merge",
		githubAPIBaseURL, c.repoOwner, c.repo, prNumber)

	mergeReq := MergeRequest{
		CommitTitle:   prTitle,
		CommitMessage: prMessage,
		MergeMethod:   "squash",
	}

	jsonData, err := json.Marshal(mergeReq)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}

	c.decorateHeaders(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	{
		if err != nil {
			return false, fmt.Errorf("failed to execute request: %w", err)
		}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return false, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	var mergeResp MergeResponse
	if err := json.Unmarshal(body, &mergeResp); err != nil {
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	return mergeResp.Merged, nil
}

func (c *GitHubClient) decorateHeaders(req *http.Request) {
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
}
