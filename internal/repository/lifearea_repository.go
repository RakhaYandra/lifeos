package repository

import "database/sql"

type LifeAreaRow struct {
	ID   int64
	Name string
}

type LifeAreaRepository struct{ DB *sql.DB }

func (r *LifeAreaRepository) List() ([]LifeAreaRow, error) {
	rows, err := r.DB.Query(`SELECT id,name FROM life_areas ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []LifeAreaRow
	for rows.Next() {
		var a LifeAreaRow
		if err := rows.Scan(&a.ID, &a.Name); err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}
