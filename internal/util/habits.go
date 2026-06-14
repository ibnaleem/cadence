package util

import (
	"database/sql"
	"encoding/json"
)

type Habit struct {
	ID          int
	Name        string
	Description string
	Frequency   string
	CreatedAt   string
} // Habit

type HabitStatus struct {
	Habit
	DoneToday bool
} // HabitStatus

func AddHabit(db *sql.DB, name, description, frequency string, embedding []float32) error {
	var embJSON []byte
	if embedding != nil {
		embJSON, _ = json.Marshal(embedding)
	} // if
	_, err := db.Exec(
		`INSERT INTO habits (name, description, frequency, embedding) VALUES (?, ?, ?, ?)`,
		name, description, frequency, embJSON,
	)
	return err
} // AddHabit

func FindSimilarHabit(db *sql.DB, embedding []float32, threshold float32) (*Habit, float32, error) {
	rows, err := db.Query(`SELECT id, name, description, frequency, created_at, embedding FROM habits WHERE embedding IS NOT NULL`)
	if err != nil {
		return nil, 0, err
	} // if
	defer rows.Close()

	var best *Habit
	var bestSim float32

	for rows.Next() {
		var h Habit
		var embJSON string
		if err := rows.Scan(&h.ID, &h.Name, &h.Description, &h.Frequency, &h.CreatedAt, &embJSON); err != nil {
			return nil, 0, err
		} // if
		var vec []float32
		if err := json.Unmarshal([]byte(embJSON), &vec); err != nil {
			continue
		} // if
		if sim := CosineSimilarity(embedding, vec); sim > bestSim {
			bestSim = sim
			h2 := h
			best = &h2
		} // if
	} // for

	if err := rows.Err(); err != nil {
		return nil, 0, err
	} // if
	if best == nil || bestSim < threshold {
		return nil, bestSim, nil
	} // if

	return best, bestSim, nil
} // FindSimilarHabit

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

func TodayStatus(db *sql.DB) ([]HabitStatus, error) {
	rows, err := db.Query(`
		SELECT id, name, description, frequency, created_at,
			EXISTS(SELECT 1 FROM habit_logs WHERE habit_id = habits.id AND logged_at = CURRENT_DATE)
		FROM habits ORDER BY id
	`)
	if err != nil {
		return nil, err
	} // if
	defer rows.Close()

	var statuses []HabitStatus

	for rows.Next() {
		var s HabitStatus
		var done int
		if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.Frequency, &s.CreatedAt, &done); err != nil {
			return nil, err
		} // if
		s.DoneToday = done > 0
		statuses = append(statuses, s)
	} // for

	return statuses, rows.Err()
} // TodayStatus

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
