package handler

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/RakhaYandra/lifeos/internal/repository"
	"github.com/RakhaYandra/lifeos/internal/service"
	"github.com/gin-gonic/gin"
)

type SubscriptionHandler struct {
	Subs *repository.SubscriptionRepository
}

type subIn struct {
	Service       string  `json:"service" binding:"required,max=200"`
	Category      string  `json:"category"`
	Cost          float64 `json:"cost" binding:"min=0"`
	Frequency     string  `json:"frequency"`
	NextBilling   string  `json:"next_billing" binding:"required"`
	PaymentMethod string  `json:"payment_method"`
	AutoRenew     *bool   `json:"auto_renew"`
	Active        *bool   `json:"active"`
	Notes         string  `json:"notes"`
}

func toSubRow(userID int64, in subIn) (*repository.SubscriptionRow, error) {
	fq := in.Frequency
	if fq == "" {
		fq = "monthly"
	}
	if !service.ValidSubFreq(fq) {
		return nil, errors.New("invalid")
	}
	if _, err := time.Parse("2006-01-02", in.NextBilling); err != nil {
		return nil, errors.New("invalid")
	}
	s := &repository.SubscriptionRow{UserID: userID, Service: in.Service, Category: in.Category,
		Cost: in.Cost, Frequency: fq, NextBilling: in.NextBilling,
		PaymentMethod: in.PaymentMethod, AutoRenew: 1, Active: 1, Notes: in.Notes}
	if in.AutoRenew != nil && !*in.AutoRenew {
		s.AutoRenew = 0
	}
	if in.Active != nil && !*in.Active {
		s.Active = 0
	}
	return s, nil
}

func subOut(s *repository.SubscriptionRow) gin.H {
	today := service.TodayWIB()
	return gin.H{"id": s.ID, "service": s.Service, "category": s.Category, "cost": s.Cost,
		"frequency": s.Frequency, "next_billing": s.NextBilling, "payment_method": s.PaymentMethod,
		"auto_renew": s.AutoRenew == 1, "active": s.Active == 1,
		"annual_cost": service.AnnualCost(s.Cost, s.Frequency),
		"days_until":  service.DaysUntil(s.NextBilling, today), "notes": s.Notes}
}

func (h *SubscriptionHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in subIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	s, err := toSubRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	id, err := h.Subs.Create(s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	s.ID = id
	c.JSON(http.StatusCreated, subOut(s))
}

func (h *SubscriptionHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Subs.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, s := range rows {
		out = append(out, subOut(s))
	}
	c.JSON(http.StatusOK, out)
}

func (h *SubscriptionHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	s, err := h.Subs.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, subOut(s))
}

func (h *SubscriptionHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Subs.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in subIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	s, err := toSubRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	s.ID = cur.ID
	s.UserID = cur.UserID
	if err := h.Subs.Update(s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, subOut(s))
}

func (h *SubscriptionHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Subs.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *SubscriptionHandler) Upcoming(c *gin.Context) {
	uid, _ := c.Get("userID")
	within, _ := strconv.Atoi(c.DefaultQuery("within_days", "14"))
	today := service.TodayWIB()
	rows, err := h.Subs.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, s := range rows {
		if s.Active != 1 {
			continue
		}
		d := service.DaysUntil(s.NextBilling, today)
		if d == nil || *d < 0 || *d > within {
			continue
		}
		out = append(out, subOut(s))
	}
	c.JSON(http.StatusOK, out)
}
