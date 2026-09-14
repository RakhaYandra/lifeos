package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/RakhaYandra/lifeos/internal/repository"
	"github.com/RakhaYandra/lifeos/internal/service"
	"github.com/gin-gonic/gin"
)

type BudgetHandler struct {
	Budgets  *repository.BudgetRepository
	Trx      *repository.TransactionRepository
	Settings *repository.SettingsRepository
}

type budgetIn struct {
	Year     int     `json:"year" binding:"required,min=2000,max=2100"`
	Month    int     `json:"month" binding:"required,min=1,max=12"`
	Category string  `json:"category" binding:"required,max=100"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
}

func (h *BudgetHandler) enrich(uid int64, b *repository.BudgetRow) gin.H {
	actual, _ := h.Trx.ExpenseSum(uid, b.Year, b.Month, b.Category)
	util := service.Utilization(actual, b.Amount)
	warn := 80.0
	if s, err := h.Settings.GetByUser(uid); err == nil {
		warn = float64(s.BudgetWarnPct)
	}
	return gin.H{"id": b.ID, "year": b.Year, "month": b.Month, "category": b.Category,
		"amount": b.Amount, "actual": actual, "remaining": b.Amount - actual,
		"util_pct": util, "status": service.BudgetStatus(util, warn)}
}

func (h *BudgetHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in budgetIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b := &repository.BudgetRow{UserID: uid.(int64), Year: in.Year, Month: in.Month, Category: in.Category, Amount: in.Amount}
	id, err := h.Budgets.Create(b)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	b.ID = id
	c.JSON(http.StatusCreated, h.enrich(uid.(int64), b))
}

func (h *BudgetHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	y, _ := strconv.Atoi(c.Query("year"))
	m, _ := strconv.Atoi(c.Query("month"))
	rows, err := h.Budgets.List(uid.(int64), y, m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, b := range rows {
		out = append(out, h.enrich(uid.(int64), b))
	}
	c.JSON(http.StatusOK, out)
}

func (h *BudgetHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	b, err := h.Budgets.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, h.enrich(uid.(int64), b))
}

func (h *BudgetHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Budgets.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in budgetIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	b := &repository.BudgetRow{ID: cur.ID, UserID: cur.UserID, Year: in.Year, Month: in.Month, Category: in.Category, Amount: in.Amount}
	if err := h.Budgets.Update(b); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, h.enrich(uid.(int64), b))
}

func (h *BudgetHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Budgets.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}
