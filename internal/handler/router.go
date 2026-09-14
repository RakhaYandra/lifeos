package handler

import (
	"net/http"

	"github.com/RakhaYandra/lifeos/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Deps struct {
	Auth     *AuthHandler
	Settings *SettingsHandler
	LifeArea *LifeAreaHandler
	Project  *ProjectHandler
	Task     *TaskHandler
	Goal     *GoalHandler
	Habit    *HabitHandler
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

	hb := auth.Group("/habits")
	hb.POST("", d.Habit.Create)
	hb.GET("", d.Habit.List)
	hb.GET("/streaks", d.Habit.Streaks)
	hb.GET("/:id", d.Habit.Get)
	hb.PUT("/:id", d.Habit.Update)
	hb.DELETE("/:id", d.Habit.Delete)
	hb.POST("/:id/log", d.Habit.Log)
	return r
}
