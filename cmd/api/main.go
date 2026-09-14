package main

import (
	"log"

	"github.com/RakhaYandra/lifeos/internal/config"
	"github.com/RakhaYandra/lifeos/internal/handler"
	"github.com/RakhaYandra/lifeos/internal/repository"
	"github.com/RakhaYandra/lifeos/internal/service"
)

func main() {
	cfg := config.Load()
	db, err := repository.Open(cfg.DBPath)
	if err != nil {
		log.Fatalf("db: %v", err)
	}
	defer db.Close() //nolint:errcheck

	users := &repository.UserRepository{DB: db}
	settings := &repository.SettingsRepository{DB: db}
	areas := &repository.LifeAreaRepository{DB: db}
	projects := &repository.ProjectRepository{DB: db}
	tasks := &repository.TaskRepository{DB: db}
	goals := &repository.GoalRepository{DB: db}
	habits := &repository.HabitRepository{DB: db}
	milestones := &repository.MilestoneRepository{DB: db}
	trx := &repository.TransactionRepository{DB: db}
	budgets := &repository.BudgetRepository{DB: db}
	subs := &repository.SubscriptionRepository{DB: db}
	hlogs := &repository.HealthLogRepository{DB: db}
	workouts := &repository.WorkoutRepository{DB: db}
	learns := &repository.LearningRepository{DB: db}
	reads := &repository.ReadingRepository{DB: db}
	reviews := &repository.ReviewRepository{DB: db}
	reminds := &repository.ReminderRepository{DB: db}
	monthly := &repository.PeriodReviewRepository{DB: db, Table: "monthly_reviews"}
	yearly := &repository.PeriodReviewRepository{DB: db, Table: "yearly_reviews", Extras: []string{"achievements", "next_year"}}
	trips := &repository.TripRepository{DB: db}
	itin := &repository.ItineraryRepository{DB: db}
	pack := &repository.PackingRepository{DB: db}
	decisions := &repository.DecisionRepository{DB: db}
	decOpts := &repository.DecisionOptionRepository{DB: db}
	decMarks := &repository.DecisionMarkRepository{DB: db}

	authH := &handler.AuthHandler{
		Svc:      &service.AuthService{Users: users, Secret: cfg.JWTSecret},
		Users:    users,
		Settings: settings,
	}
	setH := &handler.SettingsHandler{Settings: settings}
	areaH := &handler.LifeAreaHandler{Areas: areas}
	projH := &handler.ProjectHandler{Projects: projects}
	taskH := &handler.TaskHandler{Tasks: tasks}
	goalH := &handler.GoalHandler{Goals: goals}
	habitH := &handler.HabitHandler{Habits: habits}
	msH := &handler.MilestoneHandler{Milestones: milestones}
	trxH := &handler.TransactionHandler{Trx: trx}
	budH := &handler.BudgetHandler{Budgets: budgets, Trx: trx, Settings: settings}
	subH := &handler.SubscriptionHandler{Subs: subs}
	healthH := &handler.HealthHandler{Logs: hlogs, Workouts: workouts}
	learnH := &handler.LearningHandler{Learning: learns, Reading: reads}
	reviewH := &handler.ReviewHandler{Reviews: reviews, Tasks: tasks, Trx: trx, Habits: habits}
	remindH := &handler.ReminderHandler{Reminders: reminds}
	monthlyH := &handler.PeriodReviewHandler{Reviews: monthly, Tasks: tasks, Trx: trx, Habits: habits, Goals: goals, Kind: "monthly"}
	yearlyH := &handler.PeriodReviewHandler{Reviews: yearly, Tasks: tasks, Trx: trx, Habits: habits, Goals: goals, Kind: "yearly"}
	travelH := &handler.TravelHandler{Trips: trips, Itin: itin, Pack: pack}
	decisionH := &handler.DecisionHandler{Decisions: decisions, Options: decOpts, Marks: decMarks}
	dashH := &handler.DashboardHandler{Tasks: tasks, Habits: habits, Trx: trx, Goals: goals, Subs: subs, Reminds: reminds}

	r := handler.NewRouter(&handler.Deps{Auth: authH, Settings: setH, LifeArea: areaH, Project: projH, Task: taskH, Goal: goalH, Milestone: msH, Habit: habitH, Transaction: trxH, Budget: budH, Subscription: subH, Health: healthH, Learning: learnH, Review: reviewH, Reminder: remindH, Monthly: monthlyH, Yearly: yearlyH, Travel: travelH, Decision: decisionH, Dashboard: dashH}, cfg.JWTSecret, cfg.FrontendURL)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
