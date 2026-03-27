package repositories

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/hellocommit/api/internal/models"
)

type IssueRepository struct {
	db *sql.DB
}

func NewIssueRepository(db *sql.DB) *IssueRepository {
	return &IssueRepository{db: db}
}

func (r *IssueRepository) Create(issue *models.Issue) error {
	query := `INSERT OR REPLACE INTO issues (repo_id, issue_number, title, body, html_url, labels, state, author, author_url, assignee, comments, created_at, updated_at, fetched_at, is_bookmarked)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.Exec(query,
		issue.RepoID, issue.IssueNumber, issue.Title, issue.Body,
		issue.HTMLURL, issue.Labels, issue.State, issue.Author,
		issue.AuthorURL, issue.Assignee, issue.Comments, issue.CreatedAt,
		now, now, issue.IsBookmarked,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	issue.ID = id
	issue.FetchedAt = now
	return nil
}

func (r *IssueRepository) UpsertBatch(issues []*models.Issue) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`INSERT INTO issues (repo_id, issue_number, title, body, html_url, labels, state, author, author_url, assignee, comments, created_at, updated_at, fetched_at, is_bookmarked)
			  VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
			  ON CONFLICT(repo_id, issue_number) DO UPDATE SET
			    title=excluded.title, body=excluded.body, html_url=excluded.html_url,
			    labels=excluded.labels, state=excluded.state, author=excluded.author,
			    author_url=excluded.author_url, assignee=excluded.assignee,
			    comments=excluded.comments, updated_at=excluded.updated_at, fetched_at=excluded.fetched_at`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now()
	for _, issue := range issues {
		result, err := stmt.Exec(
			issue.RepoID, issue.IssueNumber, issue.Title, issue.Body,
			issue.HTMLURL, issue.Labels, issue.State, issue.Author,
			issue.AuthorURL, issue.Assignee, issue.Comments, issue.CreatedAt,
			now, now, issue.IsBookmarked,
		)
		if err != nil {
			return err
		}
		issue.ID, _ = result.LastInsertId()
		issue.FetchedAt = now
	}
	return tx.Commit()
}

// PruneClosedIssues removes issues for a repo that are no longer in the open set.
// Called after each sync, which only returns open issues from GitHub.
func (r *IssueRepository) PruneClosedIssues(repoID int64, openNumbers []int) error {
	if len(openNumbers) == 0 {
		_, err := r.db.Exec(`DELETE FROM issues WHERE repo_id = ?`, repoID)
		return err
	}
	placeholders := strings.Repeat("?,", len(openNumbers))
	placeholders = placeholders[:len(placeholders)-1]
	args := []any{repoID}
	for _, n := range openNumbers {
		args = append(args, n)
	}
	query := fmt.Sprintf(
		`DELETE FROM issues WHERE repo_id = ? AND issue_number NOT IN (%s)`,
		placeholders,
	)
	_, err := r.db.Exec(query, args...)
	return err
}

// DeleteByNumbers removes specific issues by issue number. Used on incremental syncs
// to delete issues that GitHub reported as closed since the last sync.
func (r *IssueRepository) DeleteByNumbers(repoID int64, numbers []int) error {
	if len(numbers) == 0 {
		return nil
	}
	placeholders := strings.Repeat("?,", len(numbers))
	placeholders = placeholders[:len(placeholders)-1]
	args := []any{repoID}
	for _, n := range numbers {
		args = append(args, n)
	}
	query := fmt.Sprintf(
		`DELETE FROM issues WHERE repo_id = ? AND issue_number IN (%s)`,
		placeholders,
	)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *IssueRepository) GetByID(id int64) (*models.Issue, error) {
	query := `SELECT id, repo_id, issue_number, title, body, html_url, labels, state,
			  author, author_url, assignee, comments, created_at, updated_at, fetched_at, is_bookmarked
			  FROM issues WHERE id = ?`
	issue := &models.Issue{}
	var labels, body, author, authorURL, assignee sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&issue.ID, &issue.RepoID, &issue.IssueNumber, &issue.Title,
		&body, &issue.HTMLURL, &labels, &issue.State,
		&author, &authorURL, &assignee, &issue.Comments,
		&issue.CreatedAt, &issue.UpdatedAt, &issue.FetchedAt, &issue.IsBookmarked,
	)
	if err != nil {
		return nil, err
	}
	issue.Body = body.String
	issue.Labels = labels.String
	issue.Author = author.String
	issue.AuthorURL = authorURL.String
	issue.Assignee = assignee.String
	return issue, nil
}

func (r *IssueRepository) GetByRepoID(repoID int64, labelFilter string) ([]*models.Issue, error) {
	var query string
	var args []interface{}

	if labelFilter != "" {
		query = `SELECT id, repo_id, issue_number, title, body, html_url, labels, state, author, author_url, assignee, comments, created_at, updated_at, fetched_at, is_bookmarked
				 FROM issues WHERE repo_id = ? AND labels LIKE ? AND dismissed_at IS NULL ORDER BY created_at DESC`
		args = []interface{}{repoID, "%" + labelFilter + "%"}
	} else {
		query = `SELECT id, repo_id, issue_number, title, body, html_url, labels, state, author, author_url, assignee, comments, created_at, updated_at, fetched_at, is_bookmarked
				 FROM issues WHERE repo_id = ? AND dismissed_at IS NULL ORDER BY created_at DESC`
		args = []interface{}{repoID}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanIssues(rows)
}

func (r *IssueRepository) GetByUserID(userID int64, repoID *int64, labelFilter string) ([]*models.Issue, error) {
	var query string
	var args []interface{}

	baseQuery := `SELECT i.id, i.repo_id, i.issue_number, i.title, i.body, i.html_url, i.labels, i.state, i.author, i.author_url, i.assignee, i.comments, i.created_at, i.updated_at, i.fetched_at, i.is_bookmarked
				 FROM issues i
				 INNER JOIN repos r ON i.repo_id = r.id
				 WHERE r.user_id = ? AND i.dismissed_at IS NULL`

	if repoID != nil {
		query = baseQuery + ` AND r.id = ?`
		args = []interface{}{userID, *repoID}
	} else {
		query = baseQuery
		args = []interface{}{userID}
	}

	if labelFilter != "" {
		query += ` AND i.labels LIKE ?`
		args = append(args, "%"+labelFilter+"%")
	}

	query += ` ORDER BY i.created_at DESC LIMIT 100`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanIssues(rows)
}

func (r *IssueRepository) GetGoodFirstIssues(userID int64) ([]*models.Issue, error) {
	query := `SELECT i.id, i.repo_id, i.issue_number, i.title, i.body, i.html_url, i.labels, i.state, i.author, i.author_url, i.assignee, i.comments, i.created_at, i.updated_at, i.fetched_at, i.is_bookmarked
			  FROM issues i
			  INNER JOIN repos r ON i.repo_id = r.id
			  WHERE r.user_id = ?
			  AND (i.labels LIKE '%good first issue%' OR i.labels LIKE '%good-first-issue%' OR i.labels LIKE '%beginner%' OR i.labels LIKE '%help wanted%')
			  AND i.dismissed_at IS NULL
			  ORDER BY i.created_at DESC LIMIT 100`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanIssues(rows)
}

func (r *IssueRepository) ToggleBookmark(issueID, userID int64) error {
	query := `UPDATE issues SET is_bookmarked = NOT is_bookmarked
			  WHERE id = ? AND repo_id IN (SELECT id FROM repos WHERE user_id = ?)`
	_, err := r.db.Exec(query, issueID, userID)
	return err
}

func (r *IssueRepository) Dismiss(issueID, userID int64) error {
	query := `UPDATE issues SET dismissed_at = ?
			  WHERE id = ? AND repo_id IN (SELECT id FROM repos WHERE user_id = ?)`
	_, err := r.db.Exec(query, time.Now(), issueID, userID)
	return err
}

func (r *IssueRepository) scanIssues(rows *sql.Rows) ([]*models.Issue, error) {
	var issues []*models.Issue
	for rows.Next() {
		issue := &models.Issue{}
		var labels, body, author, authorURL, assignee sql.NullString

		err := rows.Scan(
			&issue.ID, &issue.RepoID, &issue.IssueNumber, &issue.Title,
			&body, &issue.HTMLURL, &labels, &issue.State,
			&author, &authorURL, &assignee, &issue.Comments,
			&issue.CreatedAt, &issue.UpdatedAt, &issue.FetchedAt, &issue.IsBookmarked,
		)
		if err != nil {
			return nil, err
		}

		issue.Labels = labels.String
		issue.Body = body.String
		issue.Author = author.String
		issue.AuthorURL = authorURL.String
		issue.Assignee = assignee.String

		issues = append(issues, issue)
	}

	return issues, nil
}

func ParseLabels(labels []string) string {
	return strings.Join(labels, ",")
}
