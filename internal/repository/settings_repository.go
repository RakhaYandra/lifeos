package repository

import "database/sql"

type SettingsRow struct {
	UserID        int64
	ActiveYear    int
	Currency      string
	BudgetWarnPct int
	GoalWarnPct   int
}

type SettingsRepository struct{ DB *sql.DB }

func (r *SettingsRepository) GetByUser(userID int64) (*SettingsRow, error) {
	var s SettingsRow
	err := r.DB.QueryRow(`SELECT user_id,active_year,currency,budget_warn_pct,goal_warn_pct FROM settings WHERE user_id=?`, userID).
		Scan(&s.UserID, &s.ActiveYear, &s.Currency, &s.BudgetWarnPct, &s.GoalWarnPct)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *SettingsRepository) Upsert(s *SettingsRow) error {
	_, err := r.DB.Exec(`INSERT INTO settings(user_id,active_year,currency,budget_warn_pct,goal_warn_pct) VALUES(?,?,?,?,?)
		ON CONFLICT(user_id) DO UPDATE SET active_year=excluded.active_year, currency=excluded.currency,
		budget_warn_pct=excluded.budget_warn_pct, goal_warn_pct=excluded.goal_warn_pct`,
		s.UserID, s.ActiveYear, s.Currency, s.BudgetWarnPct, s.GoalWarnPct)
	return err
}

func (r *SettingsRepository) EnsureDefaults(userID int64, year int) error {
	_, err := r.DB.Exec(`INSERT OR IGNORE INTO settings(user_id,active_year,currency,budget_warn_pct,goal_warn_pct) VALUES(?,?,'IDR',80,70)`,
		userID, year)
	return err
}
