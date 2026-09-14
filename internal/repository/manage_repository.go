package repository

import "database/sql"

type SavingRow struct {
	ID                  int64
	UserID              int64
	Name                string
	TargetAmount        float64
	CurrentAmount       float64
	MonthlyContribution float64
	TargetDate          sql.NullString
	Status              string
}

type SavingRepository struct{ DB *sql.DB }

const saveCols = `id,user_id,name,target_amount,current_amount,monthly_contribution,target_date,status`

func (r *SavingRepository) Create(s *SavingRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO savings_goals(user_id,name,target_amount,current_amount,monthly_contribution,target_date,status)
		VALUES(?,?,?,?,?,?,?)`, s.UserID, s.Name, s.TargetAmount, s.CurrentAmount, s.MonthlyContribution, nullStr(s.TargetDate), s.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *SavingRepository) Get(userID, id int64) (*SavingRow, error) {
	var s SavingRow
	err := r.DB.QueryRow(`SELECT `+saveCols+` FROM savings_goals WHERE user_id=? AND id=?`, userID, id).
		Scan(&s.ID, &s.UserID, &s.Name, &s.TargetAmount, &s.CurrentAmount, &s.MonthlyContribution, &s.TargetDate, &s.Status)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SavingRepository) List(userID int64) ([]*SavingRow, error) {
	rows, err := r.DB.Query(`SELECT `+saveCols+` FROM savings_goals WHERE user_id=? ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*SavingRow
	for rows.Next() {
		var s SavingRow
		if err := rows.Scan(&s.ID, &s.UserID, &s.Name, &s.TargetAmount, &s.CurrentAmount, &s.MonthlyContribution, &s.TargetDate, &s.Status); err != nil {
			return nil, err
		}
		out = append(out, &s)
	}
	return out, rows.Err()
}

func (r *SavingRepository) Update(s *SavingRow) error {
	_, err := r.DB.Exec(`UPDATE savings_goals SET name=?,target_amount=?,current_amount=?,monthly_contribution=?,target_date=?,status=?
		WHERE user_id=? AND id=?`, s.Name, s.TargetAmount, s.CurrentAmount, s.MonthlyContribution, nullStr(s.TargetDate), s.Status, s.UserID, s.ID)
	return err
}

func (r *SavingRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM savings_goals WHERE user_id=? AND id=?`, userID, id)
	return err
}

type AssetRow struct {
	ID            int64
	UserID        int64
	Name          string
	Category      string
	PurchasePrice float64
	CurrentValue  float64
	Condition     string
	Location      string
	WarrantyEnd   sql.NullString
	Notes         string
}

type AssetRepository struct{ DB *sql.DB }

const assetCols = `id,user_id,name,category,purchase_price,current_value,condition,location,warranty_end,notes`

func (r *AssetRepository) Create(a *AssetRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO assets(user_id,name,category,purchase_price,current_value,condition,location,warranty_end,notes)
		VALUES(?,?,?,?,?,?,?,?,?)`, a.UserID, a.Name, a.Category, a.PurchasePrice, a.CurrentValue, a.Condition, a.Location, nullStr(a.WarrantyEnd), a.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *AssetRepository) List(userID int64) ([]*AssetRow, error) {
	rows, err := r.DB.Query(`SELECT `+assetCols+` FROM assets WHERE user_id=? ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*AssetRow
	for rows.Next() {
		var a AssetRow
		if err := rows.Scan(&a.ID, &a.UserID, &a.Name, &a.Category, &a.PurchasePrice, &a.CurrentValue, &a.Condition, &a.Location, &a.WarrantyEnd, &a.Notes); err != nil {
			return nil, err
		}
		out = append(out, &a)
	}
	return out, rows.Err()
}

func (r *AssetRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM assets WHERE user_id=? AND id=?`, userID, id)
	return err
}

type WishlistRow struct {
	ID         int64
	UserID     int64
	Item       string
	Category   string
	Priority   string
	EstPrice   float64
	Saved      float64
	TargetDate sql.NullString
	Status     string
}

type WishlistRepository struct{ DB *sql.DB }

const wishCols = `id,user_id,item,category,priority,est_price,saved,target_date,status`

func (r *WishlistRepository) Create(w *WishlistRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO wishlist(user_id,item,category,priority,est_price,saved,target_date,status)
		VALUES(?,?,?,?,?,?,?,?)`, w.UserID, w.Item, w.Category, w.Priority, w.EstPrice, w.Saved, nullStr(w.TargetDate), w.Status)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *WishlistRepository) List(userID int64) ([]*WishlistRow, error) {
	rows, err := r.DB.Query(`SELECT `+wishCols+` FROM wishlist WHERE user_id=? ORDER BY id`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*WishlistRow
	for rows.Next() {
		var w WishlistRow
		if err := rows.Scan(&w.ID, &w.UserID, &w.Item, &w.Category, &w.Priority, &w.EstPrice, &w.Saved, &w.TargetDate, &w.Status); err != nil {
			return nil, err
		}
		out = append(out, &w)
	}
	return out, rows.Err()
}

func (r *WishlistRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM wishlist WHERE user_id=? AND id=?`, userID, id)
	return err
}

type DocumentRow struct {
	ID           int64
	UserID       int64
	Item         string
	Category     string
	ExpiryDate   string
	ReminderDays int
	Notes        string
}

type DocumentRepository struct{ DB *sql.DB }

const docCols = `id,user_id,item,category,expiry_date,reminder_days,notes`

func (r *DocumentRepository) Create(d *DocumentRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO documents(user_id,item,category,expiry_date,reminder_days,notes)
		VALUES(?,?,?,?,?,?)`, d.UserID, d.Item, d.Category, d.ExpiryDate, d.ReminderDays, d.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *DocumentRepository) List(userID int64) ([]*DocumentRow, error) {
	rows, err := r.DB.Query(`SELECT `+docCols+` FROM documents WHERE user_id=? ORDER BY expiry_date`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*DocumentRow
	for rows.Next() {
		var d DocumentRow
		if err := rows.Scan(&d.ID, &d.UserID, &d.Item, &d.Category, &d.ExpiryDate, &d.ReminderDays, &d.Notes); err != nil {
			return nil, err
		}
		out = append(out, &d)
	}
	return out, rows.Err()
}

func (r *DocumentRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM documents WHERE user_id=? AND id=?`, userID, id)
	return err
}

type ContactRow struct {
	ID            int64
	UserID        int64
	Name          string
	Relation      string
	Birthday      sql.NullString
	LastContact   sql.NullString
	ContactMethod string
	FollowupDays  int
	Notes         string
}

type ContactRepository struct{ DB *sql.DB }

const contactCols = `id,user_id,name,relation,birthday,last_contact,contact_method,followup_days,notes`

func (r *ContactRepository) Create(x *ContactRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO contacts(user_id,name,relation,birthday,last_contact,contact_method,followup_days,notes)
		VALUES(?,?,?,?,?,?,?,?)`, x.UserID, x.Name, x.Relation, nullStr(x.Birthday), nullStr(x.LastContact), x.ContactMethod, x.FollowupDays, x.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *ContactRepository) Get(userID, id int64) (*ContactRow, error) {
	var x ContactRow
	err := r.DB.QueryRow(`SELECT `+contactCols+` FROM contacts WHERE user_id=? AND id=?`, userID, id).
		Scan(&x.ID, &x.UserID, &x.Name, &x.Relation, &x.Birthday, &x.LastContact, &x.ContactMethod, &x.FollowupDays, &x.Notes)
	if err != nil {
		return nil, err
	}
	return &x, nil
}

func (r *ContactRepository) List(userID int64) ([]*ContactRow, error) {
	rows, err := r.DB.Query(`SELECT `+contactCols+` FROM contacts WHERE user_id=? ORDER BY name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*ContactRow
	for rows.Next() {
		var x ContactRow
		if err := rows.Scan(&x.ID, &x.UserID, &x.Name, &x.Relation, &x.Birthday, &x.LastContact, &x.ContactMethod, &x.FollowupDays, &x.Notes); err != nil {
			return nil, err
		}
		out = append(out, &x)
	}
	return out, rows.Err()
}

func (r *ContactRepository) Touch(userID, id int64, date string) error {
	_, err := r.DB.Exec(`UPDATE contacts SET last_contact=? WHERE user_id=? AND id=?`, date, userID, id)
	return err
}

func (r *ContactRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM contacts WHERE user_id=? AND id=?`, userID, id)
	return err
}
