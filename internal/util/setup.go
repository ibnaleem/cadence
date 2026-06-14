package util

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func DBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	} // if

	dir := filepath.Join(home, ".cadence")

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	} // if

	return filepath.Join(dir, "cadence.db"), nil
} // DBPath

func InitDB() (*sql.DB, error) {
	path, err := DBPath()
	if err != nil {
		return nil, err
	} // if

	return sql.Open("sqlite", path)
} // InitDB

func SetupSchema(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS habits (
			id          INTEGER PRIMARY KEY,
			name        TEXT    NOT NULL,
			description TEXT,
			frequency   TEXT    NOT NULL DEFAULT 'daily',
			created_at  DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS habit_logs (
			id         INTEGER PRIMARY KEY,
			habit_id   INTEGER NOT NULL REFERENCES habits(id),
			logged_at  DATE    NOT NULL DEFAULT CURRENT_DATE,
			UNIQUE(habit_id, logged_at)
		);
	`)
	if err != nil {
		return err
	} // if

	// ponytail: ignore error — SQLite has no ADD COLUMN IF NOT EXISTS
	db.Exec(`ALTER TABLE habits ADD COLUMN embedding TEXT`)

	return nil
} // SetupSchema
