package repository

import "database/sql"

type LearningRow struct {
	ID           int64
	UserID       int64
	Topic        string
	Type         string
	Provider     string
	RelatedSkill string
	Progress     int
	Hours        float64
	Status       string
	Notes        string
}

type LearningRepository struct{ DB *sql.DB }

const learnCols = `id,user_id,topic,type,provider,related_skill,progress,hours,status,notes`

func (r *LearningRepository) Create(l *LearningRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO learning_entries(user_id,topic,type,provider,related_skill,progress,hours,status,notes)
		VALUES(?,?,?,?,?,?,?,?,?)`, l.UserID, l.Topic, l.Type, l.Provider, l.RelatedSkill, l.Progress, l.Hours, l.Status, l.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *LearningRepository) Get(userID, id int64) (*LearningRow, error) {
	var l LearningRow
	err := r.DB.QueryRow(`SELECT `+learnCols+` FROM learning_entries WHERE user_id=? AND id=?`, userID, id).
		Scan(&l.ID, &l.UserID, &l.Topic, &l.Type, &l.Provider, &l.RelatedSkill, &l.Progress, &l.Hours, &l.Status, &l.Notes)
	if err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *LearningRepository) List(userID int64) ([]*LearningRow, error) {
	rows, err := r.DB.Query(`SELECT `+learnCols+` FROM learning_entries WHERE user_id=? ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*LearningRow
	for rows.Next() {
		var l LearningRow
		if err := rows.Scan(&l.ID, &l.UserID, &l.Topic, &l.Type, &l.Provider, &l.RelatedSkill, &l.Progress, &l.Hours, &l.Status, &l.Notes); err != nil {
			return nil, err
		}
		out = append(out, &l)
	}
	return out, rows.Err()
}

func (r *LearningRepository) Update(l *LearningRow) error {
	_, err := r.DB.Exec(`UPDATE learning_entries SET topic=?,type=?,provider=?,related_skill=?,progress=?,hours=?,status=?,notes=?
		WHERE user_id=? AND id=?`, l.Topic, l.Type, l.Provider, l.RelatedSkill, l.Progress, l.Hours, l.Status, l.Notes, l.UserID, l.ID)
	return err
}

func (r *LearningRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM learning_entries WHERE user_id=? AND id=?`, userID, id)
	return err
}

type ReadingRow struct {
	ID        int64
	UserID    int64
	Title     string
	Type      string
	Author    string
	Status    string
	Rating    sql.NullInt64
	Takeaways string
}

type ReadingRepository struct{ DB *sql.DB }

const readCols = `id,user_id,title,type,author,status,rating,takeaways`

func (r *ReadingRepository) Create(x *ReadingRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO reading_entries(user_id,title,type,author,status,rating,takeaways)
		VALUES(?,?,?,?,?,?,?)`, x.UserID, x.Title, x.Type, x.Author, x.Status, nullInt64(x.Rating), x.Takeaways)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *ReadingRepository) Get(userID, id int64) (*ReadingRow, error) {
	var x ReadingRow
	err := r.DB.QueryRow(`SELECT `+readCols+` FROM reading_entries WHERE user_id=? AND id=?`, userID, id).
		Scan(&x.ID, &x.UserID, &x.Title, &x.Type, &x.Author, &x.Status, &x.Rating, &x.Takeaways)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (r *ReadingRepository) List(userID int64) ([]*ReadingRow, error) {
	rows, err := r.DB.Query(`SELECT `+readCols+` FROM reading_entries WHERE user_id=? ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*ReadingRow
	for rows.Next() {
		var x ReadingRow
		if err := rows.Scan(&x.ID, &x.UserID, &x.Title, &x.Type, &x.Author, &x.Status, &x.Rating, &x.Takeaways); err != nil {
			return nil, err
		}
		out = append(out, &x)
	}
	return out, rows.Err()
}

func (r *ReadingRepository) Update(x *ReadingRow) error {
	_, err := r.DB.Exec(`UPDATE reading_entries SET title=?,type=?,author=?,status=?,rating=?,takeaways=?
		WHERE user_id=? AND id=?`, x.Title, x.Type, x.Author, x.Status, nullInt64(x.Rating), x.Takeaways, x.UserID, x.ID)
	return err
}

func (r *ReadingRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM reading_entries WHERE user_id=? AND id=?`, userID, id)
	return err
}
