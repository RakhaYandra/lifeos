package repository

import "database/sql"

type TaskRow struct {
	ID           int64
	UserID       int64
	ProjectID    sql.NullInt64
	GoalID       sql.NullInt64
	Title        string
	Description  string
	LifeAreaID   sql.NullInt64
	Status       string
	Priority     string
	DueDate      sql.NullString
	CompletedAt  sql.NullString
	EffortEst    sql.NullFloat64
	EffortActual sql.NullFloat64
}

type TaskRepository struct{ DB *sql.DB }

const taskCols = `id,user_id,project_id,goal_id,title,description,life_area_id,status,priority,due_date,completed_at,effort_est,effort_actual`

func (r *TaskRepository) Create(t *TaskRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO tasks(user_id,project_id,goal_id,title,description,life_area_id,status,priority,due_date,completed_at,effort_est,effort_actual)
		VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`,
		t.UserID, nullInt(t.ProjectID), nullInt(t.GoalID), t.Title, t.Description, nullInt(t.LifeAreaID),
		t.Status, t.Priority, nullStr(t.DueDate), nullStr(t.CompletedAt), nullFloat(t.EffortEst), nullFloat(t.EffortActual))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *TaskRepository) Get(userID, id int64) (*TaskRow, error) {
	var t TaskRow
	err := r.DB.QueryRow(`SELECT `+taskCols+` FROM tasks WHERE user_id=? AND id=?`, userID, id).
		Scan(&t.ID, &t.UserID, &t.ProjectID, &t.GoalID, &t.Title, &t.Description, &t.LifeAreaID,
			&t.Status, &t.Priority, &t.DueDate, &t.CompletedAt, &t.EffortEst, &t.EffortActual)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

type TaskFilter struct {
	Status    string
	ProjectID int64
	DueFrom   string
	DueTo     string
}

func (r *TaskRepository) List(userID int64, f TaskFilter) ([]*TaskRow, error) {
	q := `SELECT ` + taskCols + ` FROM tasks WHERE user_id=?`
	args := []any{userID}
	if f.Status != "" {
		q += ` AND status=?`
		args = append(args, f.Status)
	}
	if f.ProjectID != 0 {
		q += ` AND project_id=?`
		args = append(args, f.ProjectID)
	}
	if f.DueFrom != "" {
		q += ` AND due_date>=?`
		args = append(args, f.DueFrom)
	}
	if f.DueTo != "" {
		q += ` AND due_date<=?`
		args = append(args, f.DueTo)
	}
	q += ` ORDER BY CASE WHEN due_date IS NULL THEN 1 ELSE 0 END, due_date, id`
	rows, err := r.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*TaskRow
	for rows.Next() {
		var t TaskRow
		if err := rows.Scan(&t.ID, &t.UserID, &t.ProjectID, &t.GoalID, &t.Title, &t.Description, &t.LifeAreaID,
			&t.Status, &t.Priority, &t.DueDate, &t.CompletedAt, &t.EffortEst, &t.EffortActual); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}

func (r *TaskRepository) Update(t *TaskRow) error {
	_, err := r.DB.Exec(`UPDATE tasks SET project_id=?,goal_id=?,title=?,description=?,life_area_id=?,status=?,priority=?,due_date=?,completed_at=?,effort_est=?,effort_actual=?
		WHERE user_id=? AND id=?`,
		nullInt(t.ProjectID), nullInt(t.GoalID), t.Title, t.Description, nullInt(t.LifeAreaID),
		t.Status, t.Priority, nullStr(t.DueDate), nullStr(t.CompletedAt), nullFloat(t.EffortEst), nullFloat(t.EffortActual),
		t.UserID, t.ID)
	return err
}

func (r *TaskRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM tasks WHERE user_id=? AND id=?`, userID, id)
	return err
}

func nullFloat(v sql.NullFloat64) any {
	if v.Valid {
		return v.Float64
	}
	return nil
}
