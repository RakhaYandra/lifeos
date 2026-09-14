package repository

import "database/sql"

type PeriodReviewRow struct {
	ID         int64
	UserID     int64
	Period     string
	Stats      string
	Wins       string
	Challenges string
	Lessons    string
	NextFocus  string
	Extras     []string
}

type PeriodReviewRepository struct {
	DB     *sql.DB
	Table  string
	Extras []string
}

func (r *PeriodReviewRepository) cols() string {
	c := `id,user_id,period,stats,wins,challenges,lessons,next_focus`
	for _, e := range r.Extras {
		c += `,` + e
	}
	return c
}

func (r *PeriodReviewRepository) placeholders() string {
	p := `?,?,?,?,?,?,?`
	for range r.Extras {
		p += `,?`
	}
	return p
}

func (r *PeriodReviewRepository) Create(v *PeriodReviewRow) (int64, error) {
	args := []any{v.UserID, v.Period, v.Stats, v.Wins, v.Challenges, v.Lessons, v.NextFocus}
	for _, e := range v.Extras {
		args = append(args, e)
	}
	res, err := r.DB.Exec(`INSERT INTO `+r.Table+`(user_id,period,stats,wins,challenges,lessons,next_focus`+
		extraSuffix(r.Extras)+`) VALUES(`+r.placeholders()+`)`, args...)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func extraSuffix(xs []string) string {
	s := ""
	for _, x := range xs {
		s += `,` + x
	}
	return s
}

func (r *PeriodReviewRepository) Get(userID, id int64) (*PeriodReviewRow, error) {
	row := r.DB.QueryRow(`SELECT `+r.cols()+` FROM `+r.Table+` WHERE user_id=? AND id=?`, userID, id)
	return r.scanRow(row)
}

func (r *PeriodReviewRepository) scanRow(row *sql.Row) (*PeriodReviewRow, error) {
	var v PeriodReviewRow
	extras := make([]sql.NullString, len(r.Extras))
	dest := []any{&v.ID, &v.UserID, &v.Period, &v.Stats, &v.Wins, &v.Challenges, &v.Lessons, &v.NextFocus}
	for i := range extras {
		dest = append(dest, &extras[i])
	}
	if err := row.Scan(dest...); err != nil {
		return nil, err
	}
	for _, e := range extras {
		if e.Valid {
			v.Extras = append(v.Extras, e.String)
		} else {
			v.Extras = append(v.Extras, "")
		}
	}
	return &v, nil
}

func (r *PeriodReviewRepository) List(userID int64) ([]*PeriodReviewRow, error) {
	rows, err := r.DB.Query(`SELECT `+r.cols()+` FROM `+r.Table+` WHERE user_id=? ORDER BY period DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*PeriodReviewRow
	for rows.Next() {
		var v PeriodReviewRow
		extras := make([]sql.NullString, len(r.Extras))
		dest := []any{&v.ID, &v.UserID, &v.Period, &v.Stats, &v.Wins, &v.Challenges, &v.Lessons, &v.NextFocus}
		for i := range extras {
			dest = append(dest, &extras[i])
		}
		if err := rows.Scan(dest...); err != nil {
			return nil, err
		}
		for _, e := range extras {
			if e.Valid {
				v.Extras = append(v.Extras, e.String)
			} else {
				v.Extras = append(v.Extras, "")
			}
		}
		out = append(out, &v)
	}
	return out, rows.Err()
}

func (r *PeriodReviewRepository) Update(v *PeriodReviewRow) error {
	set := `period=?,stats=?,wins=?,challenges=?,lessons=?,next_focus=?`
	args := []any{v.Period, v.Stats, v.Wins, v.Challenges, v.Lessons, v.NextFocus}
	for i, col := range r.Extras {
		set += `,` + col + `=?`
		if i < len(v.Extras) {
			args = append(args, v.Extras[i])
		} else {
			args = append(args, "")
		}
	}
	args = append(args, v.UserID, v.ID)
	_, err := r.DB.Exec(`UPDATE `+r.Table+` SET `+set+` WHERE user_id=? AND id=?`, args...)
	return err
}

func (r *PeriodReviewRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM `+r.Table+` WHERE user_id=? AND id=?`, userID, id)
	return err
}
