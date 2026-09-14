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

type Project struct {
	ID          int64
	UserID      int64
	Name        string
	LifeAreaID  *int64
	GoalID      *int64
	Status      string
	StartDate   string
	TargetDate  string
	CompletedAt string
	Notes       string
	Progress    int
	TotalTasks  int
	DoneTasks   int
}

type Task struct {
	ID           int64
	UserID       int64
	ProjectID    *int64
	GoalID       *int64
	Title        string
	Description  string
	LifeAreaID   *int64
	Status       string
	Priority     string
	DueDate      string
	CompletedAt  string
	EffortEst    *float64
	EffortActual *float64
	DaysLeft     *int
	Overdue      bool
}
