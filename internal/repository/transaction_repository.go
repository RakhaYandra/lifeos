package repository

import "database/sql"

type TransactionRow struct {
	ID          int64
	UserID      int64
	Date        string
	Type        string
	Category    string
	Description string
	Amount      float64
	Account     string
	Recurring   int
}

type TransactionRepository struct{ DB *sql.DB }

const trxCols = `id,user_id,date,type,category,description,amount,account,recurring`

func (r *TransactionRepository) Create(t *TransactionRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO transactions(user_id,date,type,category,description,amount,account,recurring)
		VALUES(?,?,?,?,?,?,?,?)`,
		t.UserID, t.Date, t.Type, t.Category, t.Description, t.Amount, t.Account, t.Recurring)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *TransactionRepository) Get(userID, id int64) (*TransactionRow, error) {
	var t TransactionRow
	err := r.DB.QueryRow(`SELECT `+trxCols+` FROM transactions WHERE user_id=? AND id=?`, userID, id).
		Scan(&t.ID, &t.UserID, &t.Date, &t.Type, &t.Category, &t.Description, &t.Amount, &t.Account, &t.Recurring)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

type TrxFilter struct {
	Type     string
	Category string
	Year     int
	Month    int
}

func (r *TransactionRepository) List(userID int64, f TrxFilter) ([]*TransactionRow, error) {
	q := `SELECT ` + trxCols + ` FROM transactions WHERE user_id=?`
	args := []any{userID}
	if f.Type != "" {
		q += ` AND type=?`
		args = append(args, f.Type)
	}
	if f.Category != "" {
		q += ` AND category=?`
		args = append(args, f.Category)
	}
	if f.Year != 0 {
		q += ` AND strftime('%Y',date)=?`
		args = append(args, pad4(f.Year))
	}
	if f.Month != 0 {
		q += ` AND strftime('%m',date)=?`
		args = append(args, pad2(f.Month))
	}
	q += ` ORDER BY date DESC, id DESC`
	rows, err := r.DB.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*TransactionRow
	for rows.Next() {
		var t TransactionRow
		if err := rows.Scan(&t.ID, &t.UserID, &t.Date, &t.Type, &t.Category, &t.Description, &t.Amount, &t.Account, &t.Recurring); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}

func (r *TransactionRepository) Update(t *TransactionRow) error {
	_, err := r.DB.Exec(`UPDATE transactions SET date=?,type=?,category=?,description=?,amount=?,account=?,recurring=?
		WHERE user_id=? AND id=?`,
		t.Date, t.Type, t.Category, t.Description, t.Amount, t.Account, t.Recurring, t.UserID, t.ID)
	return err
}

func (r *TransactionRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM transactions WHERE user_id=? AND id=?`, userID, id)
	return err
}

// ExpenseSum: total expense per kategori/bulan (aktual budget).
func (r *TransactionRepository) ExpenseSum(userID int64, year, month int, category string) (float64, error) {
	var s sql.NullFloat64
	err := r.DB.QueryRow(`SELECT SUM(amount) FROM transactions
		WHERE user_id=? AND type='expense' AND category=?
		AND strftime('%Y',date)=? AND strftime('%m',date)=?`,
		userID, category, pad4(year), pad2(month)).Scan(&s)
	if err != nil {
		return 0, err
	}
	if !s.Valid {
		return 0, nil
	}
	return s.Float64, nil
}

// MonthSummary: income & expense per bulan.
func (r *TransactionRepository) MonthSummary(userID int64, year, month int) (income, expense float64, err error) {
	err = r.DB.QueryRow(`SELECT COALESCE(SUM(CASE WHEN type='income' THEN amount ELSE 0 END),0),
		COALESCE(SUM(CASE WHEN type='expense' THEN amount ELSE 0 END),0) FROM transactions
		WHERE user_id=? AND strftime('%Y',date)=? AND strftime('%m',date)=?`,
		userID, pad4(year), pad2(month)).Scan(&income, &expense)
	return income, expense, err
}

// YearSummary: income & expense per tahun.
func (r *TransactionRepository) YearSummary(userID int64, year int) (income, expense float64, err error) {
	err = r.DB.QueryRow(`SELECT COALESCE(SUM(CASE WHEN type='income' THEN amount ELSE 0 END),0),
		COALESCE(SUM(CASE WHEN type='expense' THEN amount ELSE 0 END),0) FROM transactions
		WHERE user_id=? AND strftime('%Y',date)=?`,
		userID, pad4(year)).Scan(&income, &expense)
	return income, expense, err
}

func pad4(n int) string {
	if n < 10 {
		return "000" + itoa(n)
	} else if n < 100 {
		return "00" + itoa(n)
	} else if n < 1000 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

func pad2(n int) string {
	if n < 10 {
		return "0" + itoa(n)
	}
	return itoa(n)
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
