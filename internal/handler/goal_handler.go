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

type GoalHandler struct {
	Goals *repository.GoalRepository
}

type goalIn struct {
	Level        string  `json:"level"`
	ParentID     *int64  `json:"parent_id"`
	LifeAreaID   *int64  `json:"life_area_id"`
	Title        string  `json:"title" binding:"required,max=300"`
	Metric       string  `json:"metric"`
	TargetValue  float64 `json:"target_value"`
	CurrentValue float64 `json:"current_value"`
	Status       string  `json:"status"`
	TargetDate   string  `json:"target_date"`
}

func toGoalRow(userID int64, in goalIn) (*repository.GoalRow, error) {
	lv := in.Level
	if lv == "" {
		lv = "annual"
	}
	st := in.Status
	if st == "" {
		st = "active"
	}
	if !service.ValidGoalLevel(lv) || !service.ValidGoalStatus(st) {
		return nil, errors.New("invalid")
	}
	g := &repository.GoalRow{UserID: userID, Level: lv, Title: in.Title, Metric: in.Metric,
		TargetValue: in.TargetValue, CurrentValue: in.CurrentValue, Status: st}
	if in.ParentID != nil {
		g.ParentID = sql.NullInt64{Int64: *in.ParentID, Valid: true}
	}
	if in.LifeAreaID != nil {
		g.LifeAreaID = sql.NullInt64{Int64: *in.LifeAreaID, Valid: true}
	}
	if in.TargetDate != "" {
		g.TargetDate = sql.NullString{String: in.TargetDate, Valid: true}
	}
	return g, nil
}

func (h *GoalHandler) checkParent(uid int64, g *repository.GoalRow) bool {
	if !g.ParentID.Valid {
		return service.ValidGoalParent(g.Level, "")
	}
	p, err := h.Goals.Get(uid, g.ParentID.Int64)
	if err != nil {
		return false
	}
	if p.UserID != uid {
		return false
	}
	return service.ValidGoalParent(g.Level, p.Level)
}

func goalOut(g *repository.GoalRow) gin.H {
	out := gin.H{"id": g.ID, "level": g.Level, "title": g.Title, "metric": g.Metric,
		"target_value": g.TargetValue, "current_value": g.CurrentValue,
		"progress": service.GoalProgress(g.TargetValue, g.CurrentValue), "status": g.Status}
	if g.ParentID.Valid {
		out["parent_id"] = g.ParentID.Int64
	}
	if g.LifeAreaID.Valid {
		out["life_area_id"] = g.LifeAreaID.Int64
	}
	if g.TargetDate.Valid {
		out["target_date"] = g.TargetDate.String
	}
	return out
}

func (h *GoalHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in goalIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g, err := toGoalRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	if !h.checkParent(uid.(int64), g) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_parent"})
		return
	}
	id, err := h.Goals.Create(g)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	g.ID = id
	c.JSON(http.StatusCreated, goalOut(g))
}

func (h *GoalHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	lv := c.Query("level")
	if lv != "" && !service.ValidGoalLevel(lv) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_level"})
		return
	}
	rows, err := h.Goals.List(uid.(int64), lv)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, g := range rows {
		out = append(out, goalOut(g))
	}
	c.JSON(http.StatusOK, out)
}

func (h *GoalHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	g, err := h.Goals.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, goalOut(g))
}

func (h *GoalHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Goals.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in goalIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	g, err := toGoalRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	if !h.checkParent(uid.(int64), g) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_parent"})
		return
	}
	g.ID = cur.ID
	g.UserID = cur.UserID
	if err := h.Goals.Update(g); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, goalOut(g))
}

func (h *GoalHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Goals.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}
