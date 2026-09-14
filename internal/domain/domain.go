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

type Goal struct {
	ID           int64
	UserID       int64
	Level        string
	ParentID     *int64
	LifeAreaID   *int64
	Title        string
	Metric       string
	TargetValue  float64
	CurrentValue float64
	Progress     int
	Status       string
	TargetDate   string
}

type Habit struct {
	ID            int64
	UserID        int64
	Name          string
	LifeAreaID    *int64
	Frequency     string
	TargetPerWeek int
	StartDate     string
	Active        bool
	Streak        int
}

type Transaction struct {
	ID          int64
	UserID      int64
	Date        string
	Type        string
	Category    string
	Description string
	Amount      float64
	Account     string
	Recurring   bool
}

type Budget struct {
	ID        int64
	UserID    int64
	Year      int
	Month     int
	Category  string
	Amount    float64
	Actual    float64
	Remaining float64
	UtilPct   float64
	Status    string
}

type Subscription struct {
	ID            int64
	UserID        int64
	Service       string
	Category      string
	Cost          float64
	Frequency     string
	NextBilling   string
	PaymentMethod string
	AutoRenew     bool
	Active        bool
	AnnualCost    float64
	DaysUntil     *int
	Notes         string
}

type HealthLog struct {
	ID          int64
	UserID      int64
	Date        string
	Weight      *float64
	SleepHours  *float64
	WaterLiters *float64
	Energy      *int
	Mood        *int
	Notes       string
}

type Workout struct {
	ID          int64
	UserID      int64
	Date        string
	Type        string
	DurationMin int
	Intensity   string
	Calories    *int
	Notes       string
}

type LearningEntry struct {
	ID           int64
	UserID       int64
	Topic        string
	Type         string
	Provider     string
	RelatedSkill string
	Progress     int
	Hours        float64
	Status       string
	Notes        string
}

type ReadingEntry struct {
	ID        int64
	UserID    int64
	Title     string
	Type      string
	Author    string
	Status    string
	Rating    *int
	Takeaways string
}

type Review struct {
	ID         int64
	UserID     int64
	WeekStart  string
	Stats      string
	Wins       string
	Challenges string
	Lessons    string
	NextFocus  string
}

type Reminder struct {
	ID         int64
	UserID     int64
	Title      string
	Date       string
	Recurrence string
	DaysUntil  *int
	Notes      string
}
