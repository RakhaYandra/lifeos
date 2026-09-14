package repository

import "database/sql"

type GoalRow struct {
	ID           int64
	UserID       int64
	Level        string
	ParentID     sql.NullInt64
	LifeAreaID   sql.NullInt64
	Title        string
	Metric       string
	TargetValue  float64
	CurrentValue float64
	Status       string
	TargetDate   sql.NullString
}

type GoalRepository struct{ DB *sql.DB }

const goalCols = `id,user_id,level,parent_id,life_area_id,title,metric,target_value,current_value,status,target_date`

func (r *GoalRepository) Create(g *GoalRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO goals(user_id,level,parent_id,life_area_id,title,metric,target_value,current_value,status,target_date)
		VALUES(?,?,?,?,?,?,?,?,?,?)`,
		g.UserID, g.Level, nullInt(g.ParentID), nullInt(g.LifeAreaID), g.Title, g.Metric,
		g.TargetValue, g.CurrentValue, g.Status, nullStr(g.TargetDate))
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *GoalRepository) Get(userID, id int64) (*GoalRow, error) {
	var g GoalRow
	err := r.DB.QueryRow(`SELECT `+goalCols+` FROM goals WHERE user_id=? AND id=?`, userID, id).
		Scan(&g.ID, &g.UserID, &g.Level, &g.ParentID, &g.LifeAreaID, &g.Title, &g.Metric,
			&g.TargetValue, &g.CurrentValue, &g.Status, &g.TargetDate)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func (r *GoalRepository) List(userID int64, level string) ([]*GoalRow, error) {
	q := `SELECT ` + goalCols + ` FROM goals WHERE user_id=?`
	args := []any{userID}
	if level != "" {
		q += ` AND level=?`
		args = append(args, level)
	}
	q += ` ORDER BY id`
	rows, err := r.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*GoalRow
	for rows.Next() {
		var g GoalRow
		if err := rows.Scan(&g.ID, &g.UserID, &g.Level, &g.ParentID, &g.LifeAreaID, &g.Title, &g.Metric,
			&g.TargetValue, &g.CurrentValue, &g.Status, &g.TargetDate); err != nil {
			return nil, err
		}
		out = append(out, &g)
	}
	return out, rows.Err()
}

func (r *GoalRepository) Update(g *GoalRow) error {
	_, err := r.DB.Exec(`UPDATE goals SET level=?,parent_id=?,life_area_id=?,title=?,metric=?,target_value=?,current_value=?,status=?,target_date=?
		WHERE user_id=? AND id=?`,
		g.Level, nullInt(g.ParentID), nullInt(g.LifeAreaID), g.Title, g.Metric,
		g.TargetValue, g.CurrentValue, g.Status, nullStr(g.TargetDate), g.UserID, g.ID)
	return err
}

func (r *GoalRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM goals WHERE user_id=? AND id=?`, userID, id)
	return err
}
