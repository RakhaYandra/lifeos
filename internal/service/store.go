package service

import "github.com/RakhaYandra/lifeos/internal/repository"

// Seam unit test: service tergantung interface kecil, bukan *sql.DB.

type UserStore interface {
	FindByEmail(email string) (*repository.UserRow, error)
	Create(email, hash string) (int64, error)
	Count() (int, error)
}
