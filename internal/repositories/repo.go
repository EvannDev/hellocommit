package repositories

import (
	"database/sql"
	"strings"
	"time"

	"github.com/hellocommit/api/internal/models"
)

type RepoRepository struct {
	db *sql.DB
}

func NewRepoRepository(db *sql.DB) *RepoRepository {
	return &RepoRepository{db: db}
}

func (r *RepoRepository) Upsert(repo *models.Repo) error {
	query := `INSERT INTO repos (user_id, repo_id, owner, name, full_name, html_url, stars, language, fetched_at, created_at, updated_at)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			  ON CONFLICT(user_id, repo_id) DO UPDATE SET
			    name=excluded.name, full_name=excluded.full_name, html_url=excluded.html_url,
			    stars=excluded.stars, language=excluded.language, updated_at=excluded.updated_at`

	now := time.Now()
	result, err := r.db.Exec(query,
		repo.UserID, repo.RepoID, repo.Owner, repo.Name,
		repo.FullName, repo.HTMLURL, repo.Stars, repo.Language,
		now, now, now,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	repo.ID = id
	return nil
}

func (r *RepoRepository) UpdateFetchedAt(id int64, t time.Time) error {
	_, err := r.db.Exec(`UPDATE repos SET fetched_at = ? WHERE id = ?`, t, id)
	return err
}

func (r *RepoRepository) DeleteNotInList(userID int64, githubRepoIDs []int64) error {
	if len(githubRepoIDs) == 0 {
		_, err := r.db.Exec(`DELETE FROM repos WHERE user_id = ?`, userID)
		return err
	}
	placeholders := make([]string, len(githubRepoIDs))
	args := make([]interface{}, 0, len(githubRepoIDs)+1)
	args = append(args, userID)
	for i, id := range githubRepoIDs {
		placeholders[i] = "?"
		args = append(args, id)
	}
	query := "DELETE FROM repos WHERE user_id = ? AND repo_id NOT IN (" + strings.Join(placeholders, ",") + ")"
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *RepoRepository) GetByUserID(userID int64) ([]*models.Repo, error) {
	query := `SELECT id, user_id, repo_id, owner, name, full_name, html_url, stars, language, fetched_at, created_at, updated_at
			  FROM repos WHERE user_id = ? ORDER BY updated_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repos []*models.Repo
	for rows.Next() {
		repo := &models.Repo{}
		err := rows.Scan(
			&repo.ID, &repo.UserID, &repo.RepoID, &repo.Owner, &repo.Name,
			&repo.FullName, &repo.HTMLURL, &repo.Stars, &repo.Language,
			&repo.FetchedAt, &repo.CreatedAt, &repo.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		repos = append(repos, repo)
	}

	return repos, nil
}

