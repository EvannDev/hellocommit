package github

import (
	"context"
	"time"

	"github.com/google/go-github/v70/github"
)

type Client struct {
	client *github.Client
	token  string
}

func NewClient(token string) *Client {
	client := github.NewClient(nil)
	if token != "" {
		client = client.WithAuthToken(token)
	}
	return &Client{client: client, token: token}
}

func (c *Client) GetUser(ctx context.Context, username string) (*github.User, error) {
	user, _, err := c.client.Users.Get(ctx, username)
	return user, err
}

func (c *Client) GetAuthenticatedUser(ctx context.Context) (*github.User, string, error) {
	authUser, _, err := c.client.Users.Get(ctx, "")
	if err != nil {
		return nil, "", err
	}
	return authUser, c.token, nil
}

func (c *Client) ListStarredRepos(ctx context.Context, username string) ([]*github.StarredRepository, error) {
	var all []*github.StarredRepository
	opts := &github.ActivityListStarredOptions{
		Sort:        "updated",
		Direction:   "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	for {
		page, resp, err := c.client.Activity.ListStarred(ctx, username, opts)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return all, nil
}

func (c *Client) ListIssues(ctx context.Context, owner, repo string, labels []string) ([]*github.Issue, error) {
	issues, _, err := c.client.Issues.ListByRepo(ctx, owner, repo, &github.IssueListByRepoOptions{
		State:  "open",
		Labels: labels,
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	})
	return issues, err
}

func (c *Client) ListAllIssues(ctx context.Context, owner, repo string, state string, since time.Time) ([]*github.Issue, error) {
	var all []*github.Issue
	opts := &github.IssueListByRepoOptions{
		State:       state,
		ListOptions: github.ListOptions{PerPage: 100},
	}
	if !since.IsZero() {
		opts.Since = since
	}
	for {
		page, resp, err := c.client.Issues.ListByRepo(ctx, owner, repo, opts)
		if err != nil {
			return nil, err
		}
		all = append(all, page...)
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return all, nil
}

func (c *Client) GetRepo(ctx context.Context, owner, repo string) (*github.Repository, error) {
	r, _, err := c.client.Repositories.Get(ctx, owner, repo)
	return r, err
}
