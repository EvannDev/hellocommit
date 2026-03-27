package models

import "time"

type User struct {
	ID          int64     `json:"id"`
	GitHubID    int64     `json:"github_id"`
	Username    string    `json:"username"`
	AccessToken string    `json:"-"`
	LastSyncAt  time.Time `json:"last_sync_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
