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
	trx := &repository.TransactionRepository{DB: db}
	budgets := &repository.BudgetRepository{DB: db}
	subs := &repository.SubscriptionRepository{DB: db}
	hlogs := &repository.HealthLogRepository{DB: db}
	workouts := &repository.WorkoutRepository{DB: db}
	learns := &repository.LearningRepository{DB: db}
	reads := &repository.ReadingRepository{DB: db}
	reviews := &repository.ReviewRepository{DB: db}
	reminds := &repository.ReminderRepository{DB: db}

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
	trxH := &handler.TransactionHandler{Trx: trx}
	budH := &handler.BudgetHandler{Budgets: budgets, Trx: trx, Settings: settings}
	subH := &handler.SubscriptionHandler{Subs: subs}
	healthH := &handler.HealthHandler{Logs: hlogs, Workouts: workouts}
	learnH := &handler.LearningHandler{Learning: learns, Reading: reads}
	reviewH := &handler.ReviewHandler{Reviews: reviews, Tasks: tasks, Trx: trx, Habits: habits}
	remindH := &handler.ReminderHandler{Reminders: reminds}
	dashH := &handler.DashboardHandler{Tasks: tasks, Habits: habits, Trx: trx, Goals: goals, Subs: subs, Reminds: reminds}

	r := handler.NewRouter(&handler.Deps{Auth: authH, Settings: setH, LifeArea: areaH, Project: projH, Task: taskH, Goal: goalH, Habit: habitH, Transaction: trxH, Budget: budH, Subscription: subH, Health: healthH, Learning: learnH, Review: reviewH, Reminder: remindH, Dashboard: dashH}, cfg.JWTSecret, cfg.FrontendURL)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
