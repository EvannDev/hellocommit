package database

import (
	"database/sql"
	"strings"

	_ "modernc.org/sqlite"
)

func NewSQLite(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", filename)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(1)

	return db, nil
}

func Migrate(db *sql.DB) error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			github_id INTEGER UNIQUE NOT NULL,
			username TEXT NOT NULL,
			access_token TEXT NOT NULL,
			last_sync_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS repos (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			repo_id INTEGER NOT NULL,
			owner TEXT NOT NULL,
			name TEXT NOT NULL,
			full_name TEXT NOT NULL,
			html_url TEXT,
			stars INTEGER DEFAULT 0,
			language TEXT,
			fetched_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			UNIQUE(user_id, repo_id)
		)`,
		`CREATE TABLE IF NOT EXISTS issues (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			repo_id INTEGER NOT NULL,
			issue_number INTEGER NOT NULL,
			title TEXT NOT NULL,
			body TEXT,
			html_url TEXT NOT NULL,
			labels TEXT,
			state TEXT DEFAULT 'open',
			author TEXT,
			author_url TEXT,
			assignee TEXT,
			comments INTEGER DEFAULT 0,
			created_at DATETIME,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			fetched_at DATETIME,
			is_bookmarked BOOLEAN DEFAULT 0,
			FOREIGN KEY (repo_id) REFERENCES repos(id),
			UNIQUE(repo_id, issue_number)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_repos_user_id ON repos(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_issues_repo_id ON issues(repo_id)`,
		`CREATE INDEX IF NOT EXISTS idx_issues_labels ON issues(labels)`,
	}

	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}

	// Add assignee column to existing DBs (no-op if already present)
	if _, err := db.Exec(`ALTER TABLE issues ADD COLUMN assignee TEXT`); err != nil {
		if !strings.Contains(err.Error(), "duplicate column name") {
			return err
		}
	}

	// Add dismissed_at column to existing DBs (no-op if already present)
	if _, err := db.Exec(`ALTER TABLE issues ADD COLUMN dismissed_at DATETIME NULL`); err != nil {
		if !strings.Contains(err.Error(), "duplicate column name") {
			return err
		}
	}

	return nil
}
