package repository

import "database/sql"

type UserRow struct {
	ID           int64
	Email        string
	PasswordHash string
}

type UserRepository struct{ DB *sql.DB }

func (r *UserRepository) FindByEmail(email string) (*UserRow, error) {
	var u UserRow
	err := r.DB.QueryRow(`SELECT id,email,password_hash FROM users WHERE email=?`, email).
		Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) FindByID(id int64) (*UserRow, error) {
	var u UserRow
	err := r.DB.QueryRow(`SELECT id,email,password_hash FROM users WHERE id=?`, id).
		Scan(&u.ID, &u.Email, &u.PasswordHash)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *UserRepository) Count() (int, error) {
	var n int
	if err := r.DB.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&n); err != nil {
		return 0, err
	}
	return n, nil
}

func (r *UserRepository) Create(email, hash string) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO users(email,password_hash) VALUES(?,?)`, email, hash)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
