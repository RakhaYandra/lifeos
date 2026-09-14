package repository

import "database/sql"

type ProjectRow struct {
	ID          int64
	UserID      int64
	Name        string
	LifeAreaID  sql.NullInt64
	GoalID      sql.NullInt64
	Status      string
	StartDate   sql.NullString
	TargetDate  sql.NullString
	CompletedAt sql.NullString
	Notes       string
}

type ProjectRepository struct{ DB *sql.DB }

func scanProject(s func(...any) error) (*ProjectRow, error) {
	var p ProjectRow
	if err := s(&p.ID, &p.UserID, &p.Name, &p.LifeAreaID, &p.GoalID, &p.Status,
		&p.StartDate, &p.TargetDate, &p.CompletedAt, &p.Notes); err != nil {
		return nil, err
	}
	return &p, nil
}

const projectCols = `id,user_id,name,life_area_id,goal_id,status,start_date,target_date,completed_at,notes`

func (r *ProjectRepository) Create(p *ProjectRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO projects(user_id,name,life_area_id,goal_id,status,start_date,target_date,completed_at,notes)
		VALUES(?,?,?,?,?,?,?,?,?)`,
		p.UserID, p.Name, nullInt(p.LifeAreaID), nullInt(p.GoalID), p.Status,
		nullStr(p.StartDate), nullStr(p.TargetDate), nullStr(p.CompletedAt), p.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *ProjectRepository) Get(userID, id int64) (*ProjectRow, error) {
	return scanProject(r.DB.QueryRow(`SELECT `+projectCols+` FROM projects WHERE user_id=? AND id=?`, userID, id).
		Scan)
}

func (r *ProjectRepository) List(userID int64) ([]*ProjectRow, error) {
	rows, err := r.DB.Query(`SELECT `+projectCols+` FROM projects WHERE user_id=? ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*ProjectRow
	for rows.Next() {
		var p ProjectRow
		if err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.LifeAreaID, &p.GoalID, &p.Status,
			&p.StartDate, &p.TargetDate, &p.CompletedAt, &p.Notes); err != nil {
			return nil, err
		}
		out = append(out, &p)
	}
	return out, rows.Err()
}

func (r *ProjectRepository) Update(p *ProjectRow) error {
	_, err := r.DB.Exec(`UPDATE projects SET name=?,life_area_id=?,goal_id=?,status=?,start_date=?,target_date=?,completed_at=?,notes=?
		WHERE user_id=? AND id=?`,
		p.Name, nullInt(p.LifeAreaID), nullInt(p.GoalID), p.Status,
		nullStr(p.StartDate), nullStr(p.TargetDate), nullStr(p.CompletedAt), p.Notes, p.UserID, p.ID)
	return err
}

func (r *ProjectRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM projects WHERE user_id=? AND id=?`, userID, id)
	return err
}

func (r *ProjectRepository) TaskCounts(projectID int64) (total, done int, err error) {
	err = r.DB.QueryRow(`SELECT COUNT(*), COALESCE(SUM(CASE WHEN status='completed' THEN 1 ELSE 0 END),0)
		FROM tasks WHERE project_id=?`, projectID).Scan(&total, &done)
	return total, done, err
}

func nullInt(v sql.NullInt64) any {
	if v.Valid {
		return v.Int64
	}
	return nil
}

func nullStr(v sql.NullString) any {
	if v.Valid {
		return v.String
	}
	return nil
}
