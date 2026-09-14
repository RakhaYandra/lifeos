package repository

import "database/sql"

type SubscriptionRow struct {
	ID            int64
	UserID        int64
	Service       string
	Category      string
	Cost          float64
	Frequency     string
	NextBilling   string
	PaymentMethod string
	AutoRenew     int
	Active        int
	Notes         string
}

type SubscriptionRepository struct{ DB *sql.DB }

const subCols = `id,user_id,service,category,cost,frequency,next_billing,payment_method,auto_renew,active,notes`

func (r *SubscriptionRepository) Create(s *SubscriptionRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO subscriptions(user_id,service,category,cost,frequency,next_billing,payment_method,auto_renew,active,notes)
		VALUES(?,?,?,?,?,?,?,?,?,?)`,
		s.UserID, s.Service, s.Category, s.Cost, s.Frequency, s.NextBilling, s.PaymentMethod, s.AutoRenew, s.Active, s.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *SubscriptionRepository) Get(userID, id int64) (*SubscriptionRow, error) {
	var s SubscriptionRow
	err := r.DB.QueryRow(`SELECT `+subCols+` FROM subscriptions WHERE user_id=? AND id=?`, userID, id).
		Scan(&s.ID, &s.UserID, &s.Service, &s.Category, &s.Cost, &s.Frequency, &s.NextBilling,
			&s.PaymentMethod, &s.AutoRenew, &s.Active, &s.Notes)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SubscriptionRepository) List(userID int64) ([]*SubscriptionRow, error) {
	rows, err := r.DB.Query(`SELECT `+subCols+` FROM subscriptions WHERE user_id=? ORDER BY next_billing`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*SubscriptionRow
	for rows.Next() {
		var s SubscriptionRow
		if err := rows.Scan(&s.ID, &s.UserID, &s.Service, &s.Category, &s.Cost, &s.Frequency, &s.NextBilling,
			&s.PaymentMethod, &s.AutoRenew, &s.Active, &s.Notes); err != nil {
			return nil, err
		}
		out = append(out, &s)
	}
	return out, rows.Err()
}

func (r *SubscriptionRepository) Update(s *SubscriptionRow) error {
	_, err := r.DB.Exec(`UPDATE subscriptions SET service=?,category=?,cost=?,frequency=?,next_billing=?,payment_method=?,auto_renew=?,active=?,notes=?
		WHERE user_id=? AND id=?`,
		s.Service, s.Category, s.Cost, s.Frequency, s.NextBilling, s.PaymentMethod, s.AutoRenew, s.Active, s.Notes,
		s.UserID, s.ID)
	return err
}

func (r *SubscriptionRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM subscriptions WHERE user_id=? AND id=?`, userID, id)
	return err
}
