package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/RakhaYandra/lifeos/internal/repository"
	"github.com/RakhaYandra/lifeos/internal/service"
	"github.com/gin-gonic/gin"
)

type ReviewHandler struct {
	Reviews *repository.ReviewRepository
	Tasks   *repository.TaskRepository
	Trx     *repository.TransactionRepository
	Habits  *repository.HabitRepository
}

type reviewIn struct {
	WeekStart  string `json:"week_start" binding:"required"`
	Wins       string `json:"wins"`
	Challenges string `json:"challenges"`
	Lessons    string `json:"lessons"`
	NextFocus  string `json:"next_focus"`
}

func (h *ReviewHandler) buildStats(uid int64, weekStart string) string {
	s, err := time.Parse("2006-01-02", weekStart)
	if err != nil {
		return "{}"
	}
	end := s.AddDate(0, 0, 6).Format("2006-01-02")
	tasks, _ := h.Tasks.List(uid, repository.TaskFilter{DueFrom: weekStart, DueTo: end})
	done, overdue := 0, 0
	today := service.TodayWIB()
	for _, t := range tasks {
		if t.Status == "completed" {
			done++
			continue
		}
		due := ""
		if t.DueDate.Valid {
			due = t.DueDate.String
		}
		if service.IsOverdue(t.Status, due, today) {
			overdue++
		}
	}
	y, m, _ := weekMonth(weekStart)
	income, expense, _ := h.Trx.MonthSummary(uid, y, m)
	habits, _ := h.Habits.List(uid)
	habitDone := 0
	for _, hb := range habits {
		dates, _ := h.Habits.DoneDates(hb.ID)
		for _, d := range dates {
			if d >= weekStart && d <= end {
				habitDone++
			}
		}
	}
	b, _ := json.Marshal(gin.H{"tasks_due": len(tasks), "tasks_done": done, "tasks_overdue": overdue,
		"month_income": income, "month_expense": expense, "habit_checks": habitDone})
	return string(b)
}

func weekMonth(weekStart string) (int, int, error) {
	t, err := time.Parse("2006-01-02", weekStart)
	if err != nil {
		return 0, 0, err
	}
	return t.Year(), int(t.Month()), nil
}

func reviewOut(v *repository.ReviewRow) gin.H {
	var stats any
	if err := json.Unmarshal([]byte(v.Stats), &stats); err != nil {
		stats = gin.H{}
	}
	return gin.H{"id": v.ID, "week_start": v.WeekStart, "stats": stats,
		"wins": v.Wins, "challenges": v.Challenges, "lessons": v.Lessons, "next_focus": v.NextFocus}
}

func (h *ReviewHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in reviewIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	if _, err := time.Parse("2006-01-02", in.WeekStart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_date"})
		return
	}
	v := &repository.ReviewRow{UserID: uid.(int64), WeekStart: in.WeekStart,
		Wins: in.Wins, Challenges: in.Challenges, Lessons: in.Lessons, NextFocus: in.NextFocus}
	v.Stats = h.buildStats(uid.(int64), in.WeekStart)
	id, err := h.Reviews.Create(v)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	v.ID = id
	c.JSON(http.StatusCreated, reviewOut(v))
}

func (h *ReviewHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Reviews.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, v := range rows {
		out = append(out, reviewOut(v))
	}
	c.JSON(http.StatusOK, out)
}

func (h *ReviewHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	v, err := h.Reviews.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, reviewOut(v))
}

func (h *ReviewHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Reviews.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in reviewIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	if _, err := time.Parse("2006-01-02", in.WeekStart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_date"})
		return
	}
	v := &repository.ReviewRow{ID: cur.ID, UserID: cur.UserID, WeekStart: in.WeekStart,
		Wins: in.Wins, Challenges: in.Challenges, Lessons: in.Lessons, NextFocus: in.NextFocus}
	v.Stats = h.buildStats(uid.(int64), in.WeekStart)
	if err := h.Reviews.Update(v); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, reviewOut(v))
}

func (h *ReviewHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Reviews.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}
