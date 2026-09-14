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

type ReminderHandler struct {
	Reminders *repository.ReminderRepository
}

type reminderIn struct {
	Title      string `json:"title" binding:"required,max=300"`
	Date       string `json:"date" binding:"required"`
	Recurrence string `json:"recurrence"`
	Notes      string `json:"notes"`
}

func toReminderRow(userID int64, in reminderIn) (*repository.ReminderRow, error) {
	rec := in.Recurrence
	if rec == "" {
		rec = "none"
	}
	if !service.ValidRecurrence(rec) {
		return nil, errors.New("invalid")
	}
	if _, err := time.Parse("2006-01-02", in.Date); err != nil {
		return nil, errors.New("invalid")
	}
	return &repository.ReminderRow{UserID: userID, Title: in.Title, Date: in.Date, Recurrence: rec, Notes: in.Notes}, nil
}

func reminderOut(m *repository.ReminderRow) gin.H {
	today := service.TodayWIB()
	next := service.NextOccurrence(m.Date, m.Recurrence, today)
	return gin.H{"id": m.ID, "title": m.Title, "date": m.Date, "next": next,
		"recurrence": m.Recurrence, "days_until": service.DaysUntil(next, today), "notes": m.Notes}
}

func (h *ReminderHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in reminderIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m, err := toReminderRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	id, err := h.Reminders.Create(m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	m.ID = id
	c.JSON(http.StatusCreated, reminderOut(m))
}

func (h *ReminderHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Reminders.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, m := range rows {
		out = append(out, reminderOut(m))
	}
	c.JSON(http.StatusOK, out)
}

func (h *ReminderHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	m, err := h.Reminders.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, reminderOut(m))
}

func (h *ReminderHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Reminders.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in reminderIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m, err := toReminderRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	m.ID = cur.ID
	m.UserID = cur.UserID
	if err := h.Reminders.Update(m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, reminderOut(m))
}

func (h *ReminderHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Reminders.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *ReminderHandler) Upcoming(c *gin.Context) {
	uid, _ := c.Get("userID")
	within, _ := strconv.Atoi(c.DefaultQuery("within_days", "14"))
	rows, err := h.Reminders.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, m := range rows {
		o := reminderOut(m)
		d, _ := o["days_until"].(*int)
		if d == nil || *d < 0 || *d > within {
			continue
		}
		out = append(out, o)
	}
	c.JSON(http.StatusOK, out)
}
