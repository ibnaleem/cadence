package util

import "database/sql"

type Habit struct {
	ID          int
	Name        string
	Description string
	Frequency   string
	CreatedAt   string
} // Habit

func AddHabit(db *sql.DB, name, description, frequency string) error {
	_, err := db.Exec(
		`INSERT INTO habits (name, description, frequency) VALUES (?, ?, ?)`,
		name, description, frequency,
	)
	return err
} // AddHabit

func LogHabit(db *sql.DB, habitID int) error {
	_, err := db.Exec(
		`INSERT INTO habit_logs (habit_id) VALUES (?)`,
		habitID,
	)
	return err
} // LogHabit

func HabitNameByID(db *sql.DB, habitID int) (string, error) {
	var name string
	err := db.QueryRow(`SELECT name FROM habits WHERE id = ?`, habitID).Scan(&name)
	return name, err
} // HabitNameByID

func ListHabits(db *sql.DB) ([]Habit, error) {
	rows, err := db.Query(`SELECT id, name, description, frequency, created_at FROM habits ORDER BY id`)
	if err != nil {
		return nil, err
	} // if
	defer rows.Close()

	var habits []Habit

	for rows.Next() {
		var h Habit
		if err := rows.Scan(&h.ID, &h.Name, &h.Description, &h.Frequency, &h.CreatedAt); err != nil {
			return nil, err
		} // if
		habits = append(habits, h)
	} // for

	return habits, rows.Err()
} // ListHabits
