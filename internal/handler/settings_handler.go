package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/RakhaYandra/lifeos/internal/repository"
	"github.com/gin-gonic/gin"
)

type SettingsHandler struct {
	Settings *repository.SettingsRepository
}

func (h *SettingsHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	s, err := h.Settings.GetByUser(uid.(int64))
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, s)
}

type settingsIn struct {
	ActiveYear    int    `json:"active_year" binding:"required,min=2000,max=2100"`
	Currency      string `json:"currency" binding:"required,max=8"`
	BudgetWarnPct int    `json:"budget_warn_pct" binding:"min=1,max=100"`
	GoalWarnPct   int    `json:"goal_warn_pct" binding:"min=1,max=100"`
}

func (h *SettingsHandler) Put(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in settingsIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	s := &repository.SettingsRow{UserID: uid.(int64), ActiveYear: in.ActiveYear, Currency: in.Currency, BudgetWarnPct: in.BudgetWarnPct, GoalWarnPct: in.GoalWarnPct}
	if err := h.Settings.Upsert(s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, s)
}
