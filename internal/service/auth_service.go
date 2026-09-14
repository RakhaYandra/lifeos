package service

import (
	"database/sql"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("email atau password salah")
var ErrEmailTaken = errors.New("email sudah terdaftar")
var ErrSingleUser = errors.New("single_user_only")

type AuthService struct {
	Users  UserStore
	Secret string
}

func (s *AuthService) Register(email, password string) (int64, error) {
	if _, err := s.Users.FindByEmail(email); err == nil {
		return 0, ErrEmailTaken
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, err
	}
	n, err := s.Users.Count()
	if err != nil {
		return 0, err
	}
	if n >= 1 {
		return 0, ErrSingleUser
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	return s.Users.Create(email, string(hash))
}

func (s *AuthService) Login(email, password string) (string, error) {
	u, err := s.Users.FindByEmail(email)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidCredentials
	} else if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   u.ID,
		"email": u.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	})
	return tok.SignedString([]byte(s.Secret))
}
