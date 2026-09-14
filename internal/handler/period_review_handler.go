package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/RakhaYandra/lifeos/internal/repository"
	"github.com/RakhaYandra/lifeos/internal/service"
	"github.com/gin-gonic/gin"
)

type PeriodReviewHandler struct {
	Reviews *repository.PeriodReviewRepository
	Tasks   *repository.TaskRepository
	Trx     *repository.TransactionRepository
	Habits  *repository.HabitRepository
	Goals   *repository.GoalRepository
	Kind    string // "monthly" atau "yearly"
}

type periodReviewIn struct {
	Period     string `json:"period" binding:"required"`
	Wins       string `json:"wins"`
	Challenges string `json:"challenges"`
	Lessons    string `json:"lessons"`
	NextFocus  string `json:"next_focus"`
	Extra1     string `json:"extra1"`
	Extra2     string `json:"extra2"`
}

func (h *PeriodReviewHandler) validPeriod(p string) bool {
	if h.Kind == "yearly" {
		return service.ValidYearPeriod(p)
	}
	return service.ValidMonthPeriod(p)
}

func (h *PeriodReviewHandler) bounds(p string) (from, to string, ok bool) {
	if h.Kind == "yearly" {
		return service.YearBounds(p)
	}
	return service.MonthBounds(p)
}

func (h *PeriodReviewHandler) buildStats(uid int64, period string) string {
	from, to, ok := h.bounds(period)
	if !ok {
		return "{}"
	}
	tasks, _ := h.Tasks.List(uid, repository.TaskFilter{DueFrom: from, DueTo: to})
	done := 0
	for _, t := range tasks {
		if t.Status == "completed" {
			done++
		}
	}
	habits, _ := h.Habits.List(uid)
	checks := 0
	for _, hb := range habits {
		dates, _ := h.Habits.DoneDates(hb.ID)
		checks += service.CountInRange(dates, from, to)
	}
	goals, _ := h.Goals.List(uid, "")
	gDone, gActive := 0, 0
	for _, g := range goals {
		if g.Status == "completed" {
			gDone++
		} else if g.Status != "cancelled" {
			gActive++
		}
	}
	stats := gin.H{"tasks_due": len(tasks), "tasks_done": done, "habit_checks": checks,
		"goals_done": gDone, "goals_active": gActive}
	if h.Kind == "yearly" {
		y, _ := strconv.Atoi(period)
		income, expense, _ := h.Trx.YearSummary(uid, y)
		stats["year_income"] = income
		stats["year_expense"] = expense
	} else {
		var y, m int
		y, _ = strconv.Atoi(period[:4])
		m, _ = strconv.Atoi(period[5:7])
		income, expense, _ := h.Trx.MonthSummary(uid, y, m)
		stats["month_income"] = income
		stats["month_expense"] = expense
	}
	return toStatsJSON(stats)
}

func toStatsJSON(stats gin.H) string {
	b, _ := json.Marshal(stats)
	return string(b)
}

func periodOut(v *repository.PeriodReviewRow, extras []string) gin.H {
	var stats any
	if err := json.Unmarshal([]byte(v.Stats), &stats); err != nil {
		stats = gin.H{}
	}
	out := gin.H{"id": v.ID, "period": v.Period, "stats": stats,
		"wins": v.Wins, "challenges": v.Challenges, "lessons": v.Lessons, "next_focus": v.NextFocus}
	for i, name := range extras {
		if i < len(v.Extras) {
			out[name] = v.Extras[i]
		}
	}
	return out
}

func (h *PeriodReviewHandler) extraNames() []string {
	if h.Kind == "yearly" {
		return []string{"achievements", "next_year"}
	}
	return nil
}

func (h *PeriodReviewHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in periodReviewIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	if !h.validPeriod(in.Period) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_period"})
		return
	}
	v := &repository.PeriodReviewRow{UserID: uid.(int64), Period: in.Period,
		Wins: in.Wins, Challenges: in.Challenges, Lessons: in.Lessons, NextFocus: in.NextFocus}
	if h.Kind == "yearly" {
		v.Extras = []string{in.Extra1, in.Extra2}
	}
	v.Stats = h.buildStats(uid.(int64), in.Period)
	id, err := h.Reviews.Create(v)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	v.ID = id
	c.JSON(http.StatusCreated, periodOut(v, h.extraNames()))
}

func (h *PeriodReviewHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Reviews.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, v := range rows {
		out = append(out, periodOut(v, h.extraNames()))
	}
	c.JSON(http.StatusOK, out)
}

func (h *PeriodReviewHandler) Get(c *gin.Context) {
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
	c.JSON(http.StatusOK, periodOut(v, h.extraNames()))
}

func (h *PeriodReviewHandler) Update(c *gin.Context) {
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
	var in periodReviewIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	if !h.validPeriod(in.Period) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_period"})
		return
	}
	v := &repository.PeriodReviewRow{ID: cur.ID, UserID: cur.UserID, Period: in.Period,
		Wins: in.Wins, Challenges: in.Challenges, Lessons: in.Lessons, NextFocus: in.NextFocus}
	if h.Kind == "yearly" {
		v.Extras = []string{in.Extra1, in.Extra2}
	}
	v.Stats = h.buildStats(uid.(int64), in.Period)
	if err := h.Reviews.Update(v); err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "conflict"})
		return
	}
	c.JSON(http.StatusOK, periodOut(v, h.extraNames()))
}

func (h *PeriodReviewHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Reviews.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}
