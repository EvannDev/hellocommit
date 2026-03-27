package models

import "time"

type Repo struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	RepoID    int64     `json:"repo_id"`
	Owner     string    `json:"owner"`
	Name      string    `json:"name"`
	FullName  string    `json:"full_name"`
	HTMLURL   string    `json:"html_url"`
	Stars     int       `json:"stars"`
	Language  string    `json:"language"`
	FetchedAt time.Time `json:"fetched_at"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
