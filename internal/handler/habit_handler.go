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

type HabitHandler struct {
	Habits *repository.HabitRepository
}

type habitIn struct {
	Name          string `json:"name" binding:"required,max=200"`
	LifeAreaID    *int64 `json:"life_area_id"`
	Frequency     string `json:"frequency"`
	TargetPerWeek int    `json:"target_per_week"`
	StartDate     string `json:"start_date"`
	Active        *bool  `json:"active"`
}

func toHabitRow(userID int64, in habitIn) (*repository.HabitRow, error) {
	fq := in.Frequency
	if fq == "" {
		fq = "daily"
	}
	if !service.ValidHabitFreq(fq) {
		return nil, errors.New("invalid")
	}
	tpw := in.TargetPerWeek
	if tpw <= 0 {
		tpw = 1
	}
	h := &repository.HabitRow{UserID: userID, Name: in.Name, Frequency: fq, TargetPerWeek: tpw, Active: 1}
	if in.LifeAreaID != nil {
		h.LifeAreaID = sql.NullInt64{Int64: *in.LifeAreaID, Valid: true}
	}
	if in.StartDate != "" {
		if _, err := time.Parse("2006-01-02", in.StartDate); err != nil {
			return nil, errors.New("invalid")
		}
		h.StartDate = sql.NullString{String: in.StartDate, Valid: true}
	}
	if in.Active != nil && !*in.Active {
		h.Active = 0
	}
	return h, nil
}

func (h *HabitHandler) streakOf(habit *repository.HabitRow) int {
	dates, err := h.Habits.DoneDates(habit.ID)
	if err != nil {
		return 0
	}
	today := service.TodayWIB()
	if habit.Frequency == "weekly" {
		return service.WeeklyStreak(dates, habit.TargetPerWeek, today)
	}
	m := map[string]bool{}
	for _, d := range dates {
		m[d] = true
	}
	return service.DailyStreak(m, today)
}

func (h *HabitHandler) habitOut(habit *repository.HabitRow) gin.H {
	out := gin.H{"id": habit.ID, "name": habit.Name, "frequency": habit.Frequency,
		"target_per_week": habit.TargetPerWeek, "active": habit.Active == 1,
		"streak": h.streakOf(habit)}
	if habit.LifeAreaID.Valid {
		out["life_area_id"] = habit.LifeAreaID.Int64
	}
	if habit.StartDate.Valid {
		out["start_date"] = habit.StartDate.String
	}
	return out
}

func (h *HabitHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in habitIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hb, err := toHabitRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	id, err := h.Habits.Create(hb)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	hb.ID = id
	c.JSON(http.StatusCreated, h.habitOut(hb))
}

func (h *HabitHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Habits.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, hb := range rows {
		out = append(out, h.habitOut(hb))
	}
	c.JSON(http.StatusOK, out)
}

func (h *HabitHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	hb, err := h.Habits.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, h.habitOut(hb))
}

func (h *HabitHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Habits.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in habitIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hb, err := toHabitRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	hb.ID = cur.ID
	hb.UserID = cur.UserID
	if err := h.Habits.Update(hb); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, h.habitOut(hb))
}

func (h *HabitHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if _, err := h.Habits.Get(uid.(int64), id); errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	if err := h.Habits.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

type habitLogIn struct {
	Date string `json:"date"`
	Done *bool  `json:"done"`
}

func (h *HabitHandler) Log(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	hb, err := h.Habits.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in habitLogIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	date := in.Date
	if date == "" {
		date = service.TodayWIB()
	}
	if _, err := time.Parse("2006-01-02", date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_date"})
		return
	}
	done := true
	if in.Done != nil {
		done = *in.Done
	}
	if err := h.Habits.Log(hb.ID, date, done); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, h.habitOut(hb))
}

func (h *HabitHandler) Streaks(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Habits.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, hb := range rows {
		out = append(out, gin.H{"id": hb.ID, "name": hb.Name,
			"frequency": hb.Frequency, "streak": h.streakOf(hb)})
	}
	c.JSON(http.StatusOK, out)
}
