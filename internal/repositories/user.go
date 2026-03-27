package repositories

import (
	"database/sql"
	"time"

	"github.com/hellocommit/api/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `INSERT INTO users (github_id, username, access_token, last_sync_at, created_at, updated_at) 
			  VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query, user.GitHubID, user.Username, user.AccessToken, now, now, now)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

func (r *UserRepository) GetByGitHubID(githubID int64) (*models.User, error) {
	query := `SELECT id, github_id, username, access_token, last_sync_at, created_at, updated_at 
			  FROM users WHERE github_id = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, githubID).Scan(
		&user.ID, &user.GitHubID, &user.Username, &user.AccessToken,
		&user.LastSyncAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	query := `SELECT id, github_id, username, access_token, last_sync_at, created_at, updated_at 
			  FROM users WHERE id = ?`

	user := &models.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.GitHubID, &user.Username, &user.AccessToken,
		&user.LastSyncAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByAccessToken(token string) (*models.User, error) {
	query := `SELECT id, github_id, username, access_token, last_sync_at, created_at, updated_at
			  FROM users WHERE access_token = ?`
	user := &models.User{}
	err := r.db.QueryRow(query, token).Scan(
		&user.ID, &user.GitHubID, &user.Username, &user.AccessToken,
		&user.LastSyncAt, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateToken(id int64, token string) error {
	query := `UPDATE users SET access_token = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, token, time.Now(), id)
	return err
}

func (r *UserRepository) UpdateLastSync(id int64) error {
	query := `UPDATE users SET last_sync_at = ?, updated_at = ? WHERE id = ?`
	now := time.Now()
	_, err := r.db.Exec(query, now, now, id)
	return err
}

func (r *UserRepository) Delete(id int64) error {
	if _, err := r.db.Exec(`DELETE FROM issues WHERE repo_id IN (SELECT id FROM repos WHERE user_id = ?)`, id); err != nil {
		return err
	}
	if _, err := r.db.Exec(`DELETE FROM repos WHERE user_id = ?`, id); err != nil {
		return err
	}
	_, err := r.db.Exec(`DELETE FROM users WHERE id = ?`, id)
	return err
}
