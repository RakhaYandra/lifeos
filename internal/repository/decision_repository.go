package repository

import "database/sql"

type DecisionRow struct {
	ID     int64
	UserID int64
	Title  string
	Notes  string
}

type DecisionRepository struct{ DB *sql.DB }

func (r *DecisionRepository) Create(d *DecisionRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO decisions(user_id,title,notes) VALUES(?,?,?)`, d.UserID, d.Title, d.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *DecisionRepository) Get(userID, id int64) (*DecisionRow, error) {
	var d DecisionRow
	err := r.DB.QueryRow(`SELECT id,user_id,title,notes FROM decisions WHERE user_id=? AND id=?`, userID, id).
		Scan(&d.ID, &d.UserID, &d.Title, &d.Notes)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DecisionRepository) List(userID int64) ([]*DecisionRow, error) {
	rows, err := r.DB.Query(`SELECT id,user_id,title,notes FROM decisions WHERE user_id=? ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*DecisionRow
	for rows.Next() {
		var d DecisionRow
		if err := rows.Scan(&d.ID, &d.UserID, &d.Title, &d.Notes); err != nil {
			return nil, err
		}
		out = append(out, &d)
	}
	return out, rows.Err()
}

func (r *DecisionRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM decisions WHERE user_id=? AND id=?`, userID, id)
	return err
}

type DecisionOptionRow struct {
	ID         int64
	DecisionID int64
	Name       string
}

type DecisionOptionRepository struct{ DB *sql.DB }

func (r *DecisionOptionRepository) Create(decisionID int64, name string) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO decision_options(decision_id,name) VALUES(?,?)`, decisionID, name)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *DecisionOptionRepository) List(decisionID int64) ([]*DecisionOptionRow, error) {
	rows, err := r.DB.Query(`SELECT id,decision_id,name FROM decision_options WHERE decision_id=? ORDER BY id`, decisionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*DecisionOptionRow
	for rows.Next() {
		var o DecisionOptionRow
		if err := rows.Scan(&o.ID, &o.DecisionID, &o.Name); err != nil {
			return nil, err
		}
		out = append(out, &o)
	}
	return out, rows.Err()
}

func (r *DecisionOptionRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM decision_options WHERE id=?`, id)
	return err
}

type DecisionMarkRow struct {
	OptionID  int64
	Criterion string
	Weight    float64
	Score     float64
}

type DecisionMarkRepository struct{ DB *sql.DB }

func (r *DecisionMarkRepository) Upsert(optionID int64, criterion string, weight, score float64) error {
	_, err := r.DB.Exec(`INSERT INTO decision_marks(option_id,criterion,weight,score) VALUES(?,?,?,?)
		ON CONFLICT(option_id,criterion) DO UPDATE SET weight=excluded.weight, score=excluded.score`,
		optionID, criterion, weight, score)
	return err
}

func (r *DecisionMarkRepository) ListByDecision(decisionID int64) ([]*DecisionMarkRow, error) {
	rows, err := r.DB.Query(`SELECT m.option_id,m.criterion,m.weight,m.score FROM decision_marks m
		JOIN decision_options o ON o.id=m.option_id WHERE o.decision_id=?`, decisionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*DecisionMarkRow
	for rows.Next() {
		var m DecisionMarkRow
		if err := rows.Scan(&m.OptionID, &m.Criterion, &m.Weight, &m.Score); err != nil {
			return nil, err
		}
		out = append(out, &m)
	}
	return out, rows.Err()
}
