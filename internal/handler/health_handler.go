package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/RakhaYandra/lifeos/internal/repository"
	"github.com/RakhaYandra/lifeos/internal/service"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	Logs     *repository.HealthLogRepository
	Workouts *repository.WorkoutRepository
}

type healthLogIn struct {
	Date        string   `json:"date" binding:"required"`
	Weight      *float64 `json:"weight"`
	SleepHours  *float64 `json:"sleep_hours"`
	WaterLiters *float64 `json:"water_liters"`
	Energy      *int     `json:"energy"`
	Mood        *int     `json:"mood"`
	Notes       string   `json:"notes"`
}

func toHealthRow(userID int64, in healthLogIn) (*repository.HealthLogRow, error) {
	if _, err := time.Parse("2006-01-02", in.Date); err != nil {
		return nil, err
	}
	h := &repository.HealthLogRow{UserID: userID, Date: in.Date, Notes: in.Notes}
	if in.Weight != nil {
		h.Weight = sql.NullFloat64{Float64: *in.Weight, Valid: true}
	}
	if in.SleepHours != nil {
		h.SleepHours = sql.NullFloat64{Float64: *in.SleepHours, Valid: true}
	}
	if in.WaterLiters != nil {
		h.WaterLiters = sql.NullFloat64{Float64: *in.WaterLiters, Valid: true}
	}
	if in.Energy != nil {
		if *in.Energy < 1 || *in.Energy > 5 {
			return nil, errBad
		}
		h.Energy = sql.NullInt64{Int64: int64(*in.Energy), Valid: true}
	}
	if in.Mood != nil {
		if *in.Mood < 1 || *in.Mood > 5 {
			return nil, errBad
		}
		h.Mood = sql.NullInt64{Int64: int64(*in.Mood), Valid: true}
	}
	return h, nil
}

func healthOut(h *repository.HealthLogRow) gin.H {
	out := gin.H{"id": h.ID, "date": h.Date, "notes": h.Notes}
	if h.Weight.Valid {
		out["weight"] = h.Weight.Float64
	}
	if h.SleepHours.Valid {
		out["sleep_hours"] = h.SleepHours.Float64
	}
	if h.WaterLiters.Valid {
		out["water_liters"] = h.WaterLiters.Float64
	}
	if h.Energy.Valid {
		out["energy"] = h.Energy.Int64
	}
	if h.Mood.Valid {
		out["mood"] = h.Mood.Int64
	}
	return out
}

func (h *HealthHandler) PutLog(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in healthLogIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	row, err := toHealthRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	id, err := h.Logs.Upsert(row)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	row.ID = id
	c.JSON(http.StatusOK, healthOut(row))
}

func (h *HealthHandler) ListLogs(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Logs.List(uid.(int64), c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, r := range rows {
		out = append(out, healthOut(r))
	}
	c.JSON(http.StatusOK, out)
}

type workoutIn struct {
	Date        string `json:"date" binding:"required"`
	Type        string `json:"type" binding:"required,max=100"`
	DurationMin int    `json:"duration_min" binding:"min=0"`
	Intensity   string `json:"intensity"`
	Calories    *int   `json:"calories"`
	Notes       string `json:"notes"`
}

func workoutOut(w *repository.WorkoutRow) gin.H {
	out := gin.H{"id": w.ID, "date": w.Date, "type": w.Type,
		"duration_min": w.DurationMin, "intensity": w.Intensity, "notes": w.Notes}
	if w.Calories.Valid {
		out["calories"] = w.Calories.Int64
	}
	return out
}

func (h *HealthHandler) CreateWorkout(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in workoutIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := time.Parse("2006-01-02", in.Date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_date"})
		return
	}
	inT := in.Intensity
	if inT == "" {
		inT = "medium"
	}
	if !service.ValidIntensity(inT) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	w := &repository.WorkoutRow{UserID: uid.(int64), Date: in.Date, Type: in.Type,
		DurationMin: in.DurationMin, Intensity: inT, Notes: in.Notes}
	if in.Calories != nil {
		w.Calories = sql.NullInt64{Int64: int64(*in.Calories), Valid: true}
	}
	id, err := h.Workouts.Create(w)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	w.ID = id
	c.JSON(http.StatusCreated, workoutOut(w))
}

func (h *HealthHandler) ListWorkouts(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Workouts.List(uid.(int64), c.Query("from"), c.Query("to"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, w := range rows {
		out = append(out, workoutOut(w))
	}
	c.JSON(http.StatusOK, out)
}
