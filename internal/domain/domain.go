package domain

import "time"

type User struct {
	ID           int64
	Email        string
	PasswordHash string
	CreatedAt    time.Time
}

type Settings struct {
	UserID        int64
	ActiveYear    int
	Currency      string
	BudgetWarnPct int
	GoalWarnPct   int
}

type LifeArea struct {
	ID   int64
	Name string
}
