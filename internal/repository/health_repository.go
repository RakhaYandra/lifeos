package repository

import "database/sql"

type HealthLogRow struct {
	ID          int64
	UserID      int64
	Date        string
	Weight      sql.NullFloat64
	SleepHours  sql.NullFloat64
	WaterLiters sql.NullFloat64
	Energy      sql.NullInt64
	Mood        sql.NullInt64
	Notes       string
}

type HealthLogRepository struct{ DB *sql.DB }

const hlogCols = `id,user_id,date,weight,sleep_hours,water_liters,energy,mood,notes`

func (r *HealthLogRepository) Upsert(h *HealthLogRow) (int64, error) {
	var id int64
	err := r.DB.QueryRow(`SELECT id FROM health_logs WHERE user_id=? AND date=?`, h.UserID, h.Date).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	if err == sql.ErrNoRows {
		res, err := r.DB.Exec(`INSERT INTO health_logs(user_id,date,weight,sleep_hours,water_liters,energy,mood,notes)
			VALUES(?,?,?,?,?,?,?,?)`, h.UserID, h.Date, nullFloat(h.Weight), nullFloat(h.SleepHours),
			nullFloat(h.WaterLiters), nullInt64(h.Energy), nullInt64(h.Mood), h.Notes)
		if err != nil {
			return 0, err
		}
		return res.LastInsertId()
	}
	_, err = r.DB.Exec(`UPDATE health_logs SET weight=?,sleep_hours=?,water_liters=?,energy=?,mood=?,notes=?
		WHERE id=?`, nullFloat(h.Weight), nullFloat(h.SleepHours), nullFloat(h.WaterLiters),
		nullInt64(h.Energy), nullInt64(h.Mood), h.Notes, id)
	return id, err
}

func (r *HealthLogRepository) List(userID int64, from, to string) ([]*HealthLogRow, error) {
	q := `SELECT ` + hlogCols + ` FROM health_logs WHERE user_id=?`
	args := []any{userID}
	if from != "" {
		q += ` AND date>=?`
		args = append(args, from)
	}
	if to != "" {
		q += ` AND date<=?`
		args = append(args, to)
	}
	q += ` ORDER BY date DESC`
	rows, err := r.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*HealthLogRow
	for rows.Next() {
		var h HealthLogRow
		if err := rows.Scan(&h.ID, &h.UserID, &h.Date, &h.Weight, &h.SleepHours, &h.WaterLiters, &h.Energy, &h.Mood, &h.Notes); err != nil {
			return nil, err
		}
		out = append(out, &h)
	}
	return out, rows.Err()
}

func (r *HealthLogRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM health_logs WHERE user_id=? AND id=?`, userID, id)
	return err
}

func nullInt64(v sql.NullInt64) any {
	if v.Valid {
		return v.Int64
	}
	return nil
}

type WorkoutRow struct {
	ID          int64
	UserID      int64
	Date        string
	Type        string
	DurationMin int
	Intensity   string
	Calories    sql.NullInt64
	Notes       string
}

type WorkoutRepository struct{ DB *sql.DB }

const workoutCols = `id,user_id,date,type,duration_min,intensity,calories,notes`

func (r *WorkoutRepository) Create(w *WorkoutRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO workouts(user_id,date,type,duration_min,intensity,calories,notes)
		VALUES(?,?,?,?,?,?,?)`, w.UserID, w.Date, w.Type, w.DurationMin, w.Intensity, nullInt64(w.Calories), w.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *WorkoutRepository) List(userID int64, from, to string) ([]*WorkoutRow, error) {
	q := `SELECT ` + workoutCols + ` FROM workouts WHERE user_id=?`
	args := []any{userID}
	if from != "" {
		q += ` AND date>=?`
		args = append(args, from)
	}
	if to != "" {
		q += ` AND date<=?`
		args = append(args, to)
	}
	q += ` ORDER BY date DESC, id DESC`
	rows, err := r.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*WorkoutRow
	for rows.Next() {
		var w WorkoutRow
		if err := rows.Scan(&w.ID, &w.UserID, &w.Date, &w.Type, &w.DurationMin, &w.Intensity, &w.Calories, &w.Notes); err != nil {
			return nil, err
		}
		out = append(out, &w)
	}
	return out, rows.Err()
}

func (r *WorkoutRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM workouts WHERE user_id=? AND id=?`, userID, id)
	return err
}
