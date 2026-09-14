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
	return r
}
