package models

import "time"

type Issue struct {
	ID           int64     `json:"id"`
	RepoID       int64     `json:"repo_id"`
	IssueNumber  int       `json:"issue_number"`
	Title        string    `json:"title"`
	Body         string    `json:"body"`
	HTMLURL      string    `json:"html_url"`
	Labels       string    `json:"labels"`
	State        string    `json:"state"`
	Author       string    `json:"author"`
	AuthorURL    string    `json:"author_url"`
	Comments     int       `json:"comments"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	FetchedAt    time.Time `json:"fetched_at"`
	Assignee     string     `json:"assignee"`
	IsBookmarked bool       `json:"is_bookmarked"`
	DismissedAt  *time.Time `json:"dismissed_at,omitempty"`
}
