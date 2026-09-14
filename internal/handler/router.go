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
	return r
}
