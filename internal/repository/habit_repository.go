package repository

import "database/sql"

type HabitRow struct {
	ID            int64
	UserID        int64
	Name          string
	LifeAreaID    sql.NullInt64
	Frequency     string
	TargetPerWeek int
	StartDate     sql.NullString
	Active        int
}

type HabitRepository struct{ DB *sql.DB }

const habitCols = `id,user_id,name,life_area_id,frequency,target_per_week,start_date,active`

func (r *HabitRepository) Create(h *HabitRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO habits(user_id,name,life_area_id,frequency,target_per_week,start_date,active)
		VALUES(?,?,?,?,?,?,?)`,
		h.UserID, h.Name, nullInt(h.LifeAreaID), h.Frequency, h.TargetPerWeek, nullStr(h.StartDate), h.Active)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *HabitRepository) Get(userID, id int64) (*HabitRow, error) {
	var h HabitRow
	err := r.DB.QueryRow(`SELECT `+habitCols+` FROM habits WHERE user_id=? AND id=?`, userID, id).
		Scan(&h.ID, &h.UserID, &h.Name, &h.LifeAreaID, &h.Frequency, &h.TargetPerWeek, &h.StartDate, &h.Active)
	if err != nil {
		return nil, err
	}
	return &h, nil
}

func (r *HabitRepository) List(userID int64) ([]*HabitRow, error) {
	rows, err := r.DB.Query(`SELECT `+habitCols+` FROM habits WHERE user_id=? ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*HabitRow
	for rows.Next() {
		var h HabitRow
		if err := rows.Scan(&h.ID, &h.UserID, &h.Name, &h.LifeAreaID, &h.Frequency, &h.TargetPerWeek, &h.StartDate, &h.Active); err != nil {
			return nil, err
		}
		out = append(out, &h)
	}
	return out, rows.Err()
}

func (r *HabitRepository) Update(h *HabitRow) error {
	_, err := r.DB.Exec(`UPDATE habits SET name=?,life_area_id=?,frequency=?,target_per_week=?,start_date=?,active=?
		WHERE user_id=? AND id=?`,
		h.Name, nullInt(h.LifeAreaID), h.Frequency, h.TargetPerWeek, nullStr(h.StartDate), h.Active, h.UserID, h.ID)
	return err
}

func (r *HabitRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM habits WHERE user_id=? AND id=?`, userID, id)
	return err
}

func (r *HabitRepository) Log(habitID int64, date string, done bool) error {
	v := 0
	if done {
		v = 1
	}
	_, err := r.DB.Exec(`INSERT INTO habit_logs(habit_id,date,done) VALUES(?,?,?)
		ON CONFLICT(habit_id,date) DO UPDATE SET done=excluded.done`, habitID, date, v)
	return err
}

func (r *HabitRepository) DoneDates(habitID int64) ([]string, error) {
	rows, err := r.DB.Query(`SELECT date FROM habit_logs WHERE habit_id=? AND done=1 ORDER BY date`, habitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []string
	for rows.Next() {
		var d string
		if err := rows.Scan(&d); err != nil {
			return nil, err
		}
		out = append(out, d)
	}
	return out, rows.Err()
}
