package util

import "database/sql"

func AddHabit(db *sql.DB, name, description, frequency string) error {
	_, err := db.Exec(
		`INSERT INTO habits (name, description, frequency) VALUES (?, ?, ?)`,
		name, description, frequency,
	)
	return err
} // AddHabit
