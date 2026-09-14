package repository

import "database/sql"

type BudgetRow struct {
	ID       int64
	UserID   int64
	Year     int
	Month    int
	Category string
	Amount   float64
}

type BudgetRepository struct{ DB *sql.DB }

const budgetCols = `id,user_id,year,month,category,amount`

func (r *BudgetRepository) Create(b *BudgetRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO budgets(user_id,year,month,category,amount) VALUES(?,?,?,?,?)`,
		b.UserID, b.Year, b.Month, b.Category, b.Amount)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *BudgetRepository) Get(userID, id int64) (*BudgetRow, error) {
	var b BudgetRow
	err := r.DB.QueryRow(`SELECT `+budgetCols+` FROM budgets WHERE user_id=? AND id=?`, userID, id).
		Scan(&b.ID, &b.UserID, &b.Year, &b.Month, &b.Category, &b.Amount)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BudgetRepository) List(userID int64, year, month int) ([]*BudgetRow, error) {
	q := `SELECT ` + budgetCols + ` FROM budgets WHERE user_id=?`
	args := []any{userID}
	if year != 0 {
		q += ` AND year=?`
		args = append(args, year)
	}
	if month != 0 {
		q += ` AND month=?`
		args = append(args, month)
	}
	q += ` ORDER BY year,month,category`
	rows, err := r.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*BudgetRow
	for rows.Next() {
		var b BudgetRow
		if err := rows.Scan(&b.ID, &b.UserID, &b.Year, &b.Month, &b.Category, &b.Amount); err != nil {
			return nil, err
		}
		out = append(out, &b)
	}
	return out, rows.Err()
}

func (r *BudgetRepository) Update(b *BudgetRow) error {
	_, err := r.DB.Exec(`UPDATE budgets SET year=?,month=?,category=?,amount=? WHERE user_id=? AND id=?`,
		b.Year, b.Month, b.Category, b.Amount, b.UserID, b.ID)
	return err
}

func (r *BudgetRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM budgets WHERE user_id=? AND id=?`, userID, id)
	return err
}
