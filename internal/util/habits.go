package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
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

func GetHabit(db *sql.DB, id int) (*Habit, error) {
	var h Habit
	err := db.QueryRow(`SELECT id, name, description, frequency, created_at FROM habits WHERE id = ?`, id).
		Scan(&h.ID, &h.Name, &h.Description, &h.Frequency, &h.CreatedAt)
	if err != nil {
		return nil, err
	} // if
	return &h, nil
} // GetHabit

func UpdateHabit(db *sql.DB, id int, name, description string, embedding []float32) error {
	if embedding != nil {
		var embJSON []byte
		embJSON, _ = json.Marshal(embedding)
		_, err := db.Exec(
			`UPDATE habits SET name = ?, description = ?, embedding = ? WHERE id = ?`,
			name, description, embJSON, id,
		)
		return err
	}
	_, err := db.Exec(`UPDATE habits SET name = ?, description = ? WHERE id = ?`, name, description, id)
	return err
} // UpdateHabit

func DeleteHabit(db *sql.DB, habitID int) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM habit_logs WHERE habit_id = ?`, habitID); err != nil {
		tx.Rollback()
		return err
	}
	if _, err := tx.Exec(`DELETE FROM habits WHERE id = ?`, habitID); err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
} // DeleteHabit

func LogHabit(db *sql.DB, habitID int, date string) (bool, error) {
	var res sql.Result
	var err error
	if date == "" {
		res, err = db.Exec(
			`INSERT OR IGNORE INTO habit_logs (habit_id) VALUES (?)`,
			habitID,
		)
	} else {
		res, err = db.Exec(
			`INSERT OR IGNORE INTO habit_logs (habit_id, logged_at) VALUES (?, ?)`,
			habitID, date,
		)
	}
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("reading rows affected: %w", err)
	}
	return n > 0, nil
} // LogHabit

func UnlogHabit(db *sql.DB, habitID int) (bool, error) {
	res, err := db.Exec(
		`DELETE FROM habit_logs WHERE habit_id = ? AND logged_at = CURRENT_DATE`,
		habitID,
	)
	if err != nil {
		return false, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("reading rows affected: %w", err)
	}
	return n > 0, nil
} // UnlogHabit

func HabitNameByID(db *sql.DB, habitID int) (string, error) {
	var name string
	err := db.QueryRow(`SELECT name FROM habits WHERE id = ?`, habitID).Scan(&name)
	return name, err
} // HabitNameByID

func TodayStatus(db *sql.DB) ([]HabitStatus, error) {
	rows, err := db.Query(`
		SELECT id, name, description, frequency, created_at,
			CASE
				WHEN frequency = 'weekly' THEN
					EXISTS(SELECT 1 FROM habit_logs WHERE habit_id = habits.id AND logged_at >= DATE('now', '-6 days'))
				ELSE
					EXISTS(SELECT 1 FROM habit_logs WHERE habit_id = habits.id AND logged_at = CURRENT_DATE)
			END
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

func WeeklyLogs(db *sql.DB) (map[int]int, error) {
	rows, err := db.Query(`
		SELECT habit_id, COUNT(*) FROM habit_logs
		WHERE logged_at >= DATE('now', '-6 days')
		  AND logged_at <= DATE('now')
		GROUP BY habit_id
	`)
	if err != nil {
		return nil, err
	} // if
	defer rows.Close()

	counts := make(map[int]int)
	for rows.Next() {
		var id, count int
		if err := rows.Scan(&id, &count); err != nil {
			return nil, err
		} // if
		counts[id] = count
	} // for

	return counts, rows.Err()
} // WeeklyLogs

func AllStreaks(db *sql.DB) (map[int]int, error) {
	rows, err := db.Query(`
		SELECT habit_id, date(logged_at) FROM habit_logs
		GROUP BY habit_id, date(logged_at)
		ORDER BY habit_id, date(logged_at) DESC
	`)
	if err != nil {
		return nil, err
	} // if
	defer rows.Close()

	byHabit := make(map[int][]string)
	for rows.Next() {
		var id int
		var date string
		if err := rows.Scan(&id, &date); err != nil {
			return nil, err
		} // if
		byHabit[id] = append(byHabit[id], date)
	} // for
	if err := rows.Err(); err != nil {
		return nil, err
	} // if

	streaks := make(map[int]int)
	for id, dates := range byHabit {
		streaks[id] = calcStreak(dates)
	} // for
	return streaks, nil
} // AllStreaks

func calcStreak(dates []string) int {
	if len(dates) == 0 {
		return 0
	} // if
	now := time.Now()
	today := now.Format("2006-01-02")
	yesterday := now.AddDate(0, 0, -1).Format("2006-01-02")
	if dates[0] != today && dates[0] != yesterday {
		return 0
	} // if
	streak := 1
	for i := 1; i < len(dates); i++ {
		prev, _ := time.Parse("2006-01-02", dates[i-1])
		curr, _ := time.Parse("2006-01-02", dates[i])
		if prev.AddDate(0, 0, -1).Equal(curr) {
			streak++
		} else {
			break
		} // if
	} // for
	return streak
} // calcStreak

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
