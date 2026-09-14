package handler

import (
	"net/http"

	"github.com/RakhaYandra/lifeos/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	Auth         *AuthHandler
	Settings     *SettingsHandler
	LifeArea     *LifeAreaHandler
	Project      *ProjectHandler
	Task         *TaskHandler
	Goal         *GoalHandler
	Milestone    *MilestoneHandler
	Habit        *HabitHandler
	Transaction  *TransactionHandler
	Budget       *BudgetHandler
	Subscription *SubscriptionHandler
	Health       *HealthHandler
	Learning     *LearningHandler
	Review       *ReviewHandler
	Reminder     *ReminderHandler
	Dashboard    *DashboardHandler
}

func NewRouter(d *Deps, secret, frontendURL string) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery(), middleware.CORS(frontendURL))
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })
	v1 := r.Group("/v1")
	v1.POST("/auth/register", d.Auth.Register)
	v1.POST("/auth/login", d.Auth.Login)

	auth := v1.Group("", middleware.Auth(secret))
	auth.GET("/me", d.Auth.Me)
	auth.GET("/settings", d.Settings.Get)
	auth.PUT("/settings", d.Settings.Put)
	auth.GET("/life-areas", d.LifeArea.List)

	pr := auth.Group("/projects")
	pr.POST("", d.Project.Create)
	pr.GET("", d.Project.List)
	pr.GET("/:id", d.Project.Get)
	pr.PUT("/:id", d.Project.Update)
	pr.DELETE("/:id", d.Project.Delete)

	tk := auth.Group("/tasks")
	tk.POST("", d.Task.Create)
	tk.GET("", d.Task.List)
	tk.GET("/today", d.Task.Today)
	tk.GET("/week", d.Task.Week)
	tk.GET("/:id", d.Task.Get)
	tk.PUT("/:id", d.Task.Update)
	tk.DELETE("/:id", d.Task.Delete)

	gl := auth.Group("/goals")
	gl.POST("", d.Goal.Create)
	gl.GET("", d.Goal.List)
	gl.GET("/:id", d.Goal.Get)
	gl.PUT("/:id", d.Goal.Update)
	gl.DELETE("/:id", d.Goal.Delete)

	ms := auth.Group("/milestones")
	ms.POST("", d.Milestone.Create)
	ms.GET("", d.Milestone.List)
	ms.GET("/upcoming", d.Milestone.Upcoming)
	ms.GET("/:id", d.Milestone.Get)
	ms.PUT("/:id", d.Milestone.Update)
	ms.DELETE("/:id", d.Milestone.Delete)

	hb := auth.Group("/habits")
	hb.POST("", d.Habit.Create)
	hb.GET("", d.Habit.List)
	hb.GET("/streaks", d.Habit.Streaks)
	hb.GET("/:id", d.Habit.Get)
	hb.PUT("/:id", d.Habit.Update)
	hb.DELETE("/:id", d.Habit.Delete)
	hb.POST("/:id/log", d.Habit.Log)
	hb.GET("/:id/logs", d.Habit.Logs)

	tx := auth.Group("/transactions")
	tx.POST("", d.Transaction.Create)
	tx.GET("", d.Transaction.List)
	tx.GET("/summary", d.Transaction.Summary)
	tx.GET("/:id", d.Transaction.Get)
	tx.PUT("/:id", d.Transaction.Update)
	tx.DELETE("/:id", d.Transaction.Delete)

	bd := auth.Group("/budgets")
	bd.POST("", d.Budget.Create)
	bd.GET("", d.Budget.List)
	bd.GET("/:id", d.Budget.Get)
	bd.PUT("/:id", d.Budget.Update)
	bd.DELETE("/:id", d.Budget.Delete)

	sb := auth.Group("/subscriptions")
	sb.POST("", d.Subscription.Create)
	sb.GET("", d.Subscription.List)
	sb.GET("/upcoming", d.Subscription.Upcoming)
	sb.GET("/:id", d.Subscription.Get)
	sb.PUT("/:id", d.Subscription.Update)
	sb.DELETE("/:id", d.Subscription.Delete)

	auth.GET("/dashboard", d.Dashboard.Get)

	hl := auth.Group("/health-logs")
	hl.PUT("", d.Health.PutLog)
	hl.GET("", d.Health.ListLogs)

	wo := auth.Group("/workouts")
	wo.POST("", d.Health.CreateWorkout)
	wo.GET("", d.Health.ListWorkouts)

	ln := auth.Group("/learning")
	ln.POST("", d.Learning.CreateLearn)
	ln.GET("", d.Learning.ListLearn)
	ln.PUT("/:id", d.Learning.UpdateLearn)
	ln.DELETE("/:id", d.Learning.DeleteLearn)

	rd := auth.Group("/reading")
	rd.POST("", d.Learning.CreateRead)
	rd.GET("", d.Learning.ListRead)
	rd.PUT("/:id", d.Learning.UpdateRead)
	rd.DELETE("/:id", d.Learning.DeleteRead)

	rv := auth.Group("/reviews")
	rv.POST("", d.Review.Create)
	rv.GET("", d.Review.List)
	rv.GET("/:id", d.Review.Get)
	rv.PUT("/:id", d.Review.Update)
	rv.DELETE("/:id", d.Review.Delete)

	rm := auth.Group("/reminders")
	rm.POST("", d.Reminder.Create)
	rm.GET("", d.Reminder.List)
	rm.GET("/upcoming", d.Reminder.Upcoming)
	rm.GET("/:id", d.Reminder.Get)
	rm.PUT("/:id", d.Reminder.Update)
	rm.DELETE("/:id", d.Reminder.Delete)
	return r
}
