package repository

import "database/sql"

type MilestoneRow struct {
	ID          int64
	UserID      int64
	GoalID      sql.NullInt64
	ProjectID   sql.NullInt64
	Title       string
	TargetDate  string
	Status      string
	CompletedAt sql.NullString
	Notes       string
}

type MilestoneRepository struct{ DB *sql.DB }

const msCols = `id,user_id,goal_id,project_id,title,target_date,status,completed_at,notes`

func (r *MilestoneRepository) Create(m *MilestoneRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO milestones(user_id,goal_id,project_id,title,target_date,status,completed_at,notes)
		VALUES(?,?,?,?,?,?,?,?)`,
		m.UserID, nullInt(m.GoalID), nullInt(m.ProjectID), m.Title, m.TargetDate, m.Status, nullStr(m.CompletedAt), m.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *MilestoneRepository) Get(userID, id int64) (*MilestoneRow, error) {
	var m MilestoneRow
	err := r.DB.QueryRow(`SELECT `+msCols+` FROM milestones WHERE user_id=? AND id=?`, userID, id).
		Scan(&m.ID, &m.UserID, &m.GoalID, &m.ProjectID, &m.Title, &m.TargetDate, &m.Status, &m.CompletedAt, &m.Notes)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

type MilestoneFilter struct {
	GoalID    int64
	ProjectID int64
}

func (r *MilestoneRepository) List(userID int64, f MilestoneFilter) ([]*MilestoneRow, error) {
	q := `SELECT ` + msCols + ` FROM milestones WHERE user_id=?`
	args := []any{userID}
	if f.GoalID != 0 {
		q += ` AND goal_id=?`
		args = append(args, f.GoalID)
	}
	if f.ProjectID != 0 {
		q += ` AND project_id=?`
		args = append(args, f.ProjectID)
	}
	q += ` ORDER BY target_date, id`
	rows, err := r.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*MilestoneRow
	for rows.Next() {
		var m MilestoneRow
		if err := rows.Scan(&m.ID, &m.UserID, &m.GoalID, &m.ProjectID, &m.Title, &m.TargetDate, &m.Status, &m.CompletedAt, &m.Notes); err != nil {
			return nil, err
		}
		out = append(out, &m)
	}
	return out, rows.Err()
}

func (r *MilestoneRepository) Update(m *MilestoneRow) error {
	_, err := r.DB.Exec(`UPDATE milestones SET goal_id=?,project_id=?,title=?,target_date=?,status=?,completed_at=?,notes=?
		WHERE user_id=? AND id=?`,
		nullInt(m.GoalID), nullInt(m.ProjectID), m.Title, m.TargetDate, m.Status, nullStr(m.CompletedAt), m.Notes,
		m.UserID, m.ID)
	return err
}

func (r *MilestoneRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM milestones WHERE user_id=? AND id=?`, userID, id)
	return err
}
