package repository

import "database/sql"

type ReviewRow struct {
	ID         int64
	UserID     int64
	WeekStart  string
	Stats      string
	Wins       string
	Challenges string
	Lessons    string
	NextFocus  string
}

type ReviewRepository struct{ DB *sql.DB }

const reviewCols = `id,user_id,week_start,stats,wins,challenges,lessons,next_focus`

func (r *ReviewRepository) Create(v *ReviewRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO reviews(user_id,week_start,stats,wins,challenges,lessons,next_focus)
		VALUES(?,?,?,?,?,?,?)`, v.UserID, v.WeekStart, v.Stats, v.Wins, v.Challenges, v.Lessons, v.NextFocus)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *ReviewRepository) Get(userID, id int64) (*ReviewRow, error) {
	var v ReviewRow
	err := r.DB.QueryRow(`SELECT `+reviewCols+` FROM reviews WHERE user_id=? AND id=?`, userID, id).
		Scan(&v.ID, &v.UserID, &v.WeekStart, &v.Stats, &v.Wins, &v.Challenges, &v.Lessons, &v.NextFocus)
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *ReviewRepository) List(userID int64) ([]*ReviewRow, error) {
	rows, err := r.DB.Query(`SELECT `+reviewCols+` FROM reviews WHERE user_id=? ORDER BY week_start DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*ReviewRow
	for rows.Next() {
		var v ReviewRow
		if err := rows.Scan(&v.ID, &v.UserID, &v.WeekStart, &v.Stats, &v.Wins, &v.Challenges, &v.Lessons, &v.NextFocus); err != nil {
			return nil, err
		}
		out = append(out, &v)
	}
	return out, rows.Err()
}

func (r *ReviewRepository) Update(v *ReviewRow) error {
	_, err := r.DB.Exec(`UPDATE reviews SET week_start=?,stats=?,wins=?,challenges=?,lessons=?,next_focus=?
		WHERE user_id=? AND id=?`, v.WeekStart, v.Stats, v.Wins, v.Challenges, v.Lessons, v.NextFocus, v.UserID, v.ID)
	return err
}

func (r *ReviewRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM reviews WHERE user_id=? AND id=?`, userID, id)
	return err
}

type ReminderRow struct {
	ID         int64
	UserID     int64
	Title      string
	Date       string
	Recurrence string
	Notes      string
}

type ReminderRepository struct{ DB *sql.DB }

const reminderCols = `id,user_id,title,date,recurrence,notes`

func (r *ReminderRepository) Create(m *ReminderRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO reminders(user_id,title,date,recurrence,notes) VALUES(?,?,?,?,?)`,
		m.UserID, m.Title, m.Date, m.Recurrence, m.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *ReminderRepository) Get(userID, id int64) (*ReminderRow, error) {
	var m ReminderRow
	err := r.DB.QueryRow(`SELECT `+reminderCols+` FROM reminders WHERE user_id=? AND id=?`, userID, id).
		Scan(&m.ID, &m.UserID, &m.Title, &m.Date, &m.Recurrence, &m.Notes)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *ReminderRepository) List(userID int64) ([]*ReminderRow, error) {
	rows, err := r.DB.Query(`SELECT `+reminderCols+` FROM reminders WHERE user_id=? ORDER BY date`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*ReminderRow
	for rows.Next() {
		var m ReminderRow
		if err := rows.Scan(&m.ID, &m.UserID, &m.Title, &m.Date, &m.Recurrence, &m.Notes); err != nil {
			return nil, err
		}
		out = append(out, &m)
	}
	return out, rows.Err()
}

func (r *ReminderRepository) Update(m *ReminderRow) error {
	_, err := r.DB.Exec(`UPDATE reminders SET title=?,date=?,recurrence=?,notes=? WHERE user_id=? AND id=?`,
		m.Title, m.Date, m.Recurrence, m.Notes, m.UserID, m.ID)
	return err
}

func (r *ReminderRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM reminders WHERE user_id=? AND id=?`, userID, id)
	return err
}
