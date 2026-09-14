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

type MilestoneHandler struct {
	Milestones *repository.MilestoneRepository
}

type milestoneIn struct {
	GoalID     *int64 `json:"goal_id"`
	ProjectID  *int64 `json:"project_id"`
	Title      string `json:"title" binding:"required,max=300"`
	TargetDate string `json:"target_date" binding:"required"`
	Status     string `json:"status"`
	Notes      string `json:"notes"`
}

func toMilestoneRow(userID int64, in milestoneIn) (*repository.MilestoneRow, error) {
	st := in.Status
	if st == "" {
		st = "pending"
	}
	if st != "pending" && st != "completed" && st != "cancelled" {
		return nil, errors.New("invalid")
	}
	if _, err := time.Parse("2006-01-02", in.TargetDate); err != nil {
		return nil, errors.New("invalid")
	}
	m := &repository.MilestoneRow{UserID: userID, Title: in.Title, TargetDate: in.TargetDate, Status: st, Notes: in.Notes}
	if in.GoalID != nil {
		m.GoalID = sql.NullInt64{Int64: *in.GoalID, Valid: true}
	}
	if in.ProjectID != nil {
		m.ProjectID = sql.NullInt64{Int64: *in.ProjectID, Valid: true}
	}
	if st == "completed" {
		m.CompletedAt = sql.NullString{String: service.TodayWIB(), Valid: true}
	}
	return m, nil
}

func milestoneOut(m *repository.MilestoneRow) gin.H {
	out := gin.H{"id": m.ID, "title": m.Title, "target_date": m.TargetDate, "status": m.Status,
		"overdue": service.IsOverdue(m.Status, m.TargetDate, service.TodayWIB()), "notes": m.Notes}
	if m.GoalID.Valid {
		out["goal_id"] = m.GoalID.Int64
	}
	if m.ProjectID.Valid {
		out["project_id"] = m.ProjectID.Int64
	}
	if m.CompletedAt.Valid {
		out["completed_at"] = m.CompletedAt.String
	}
	return out
}

func (h *MilestoneHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in milestoneIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	m, err := toMilestoneRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	id, err := h.Milestones.Create(m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	m.ID = id
	c.JSON(http.StatusCreated, milestoneOut(m))
}

func (h *MilestoneHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	var f repository.MilestoneFilter
	if g := c.Query("goal_id"); g != "" {
		f.GoalID, _ = strconv.ParseInt(g, 10, 64)
	}
	if p := c.Query("project_id"); p != "" {
		f.ProjectID, _ = strconv.ParseInt(p, 10, 64)
	}
	rows, err := h.Milestones.List(uid.(int64), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, m := range rows {
		out = append(out, milestoneOut(m))
	}
	c.JSON(http.StatusOK, out)
}

func (h *MilestoneHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	m, err := h.Milestones.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, milestoneOut(m))
}

func (h *MilestoneHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Milestones.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in milestoneIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	m, err := toMilestoneRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	m.ID = cur.ID
	m.UserID = cur.UserID
	if m.Status != "completed" {
		m.CompletedAt = sql.NullString{}
	} else if !cur.CompletedAt.Valid {
		m.CompletedAt = sql.NullString{String: service.TodayWIB(), Valid: true}
	} else {
		m.CompletedAt = cur.CompletedAt
	}
	if err := h.Milestones.Update(m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, milestoneOut(m))
}

func (h *MilestoneHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Milestones.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *MilestoneHandler) Upcoming(c *gin.Context) {
	uid, _ := c.Get("userID")
	within, _ := strconv.Atoi(c.DefaultQuery("within_days", "14"))
	today := service.TodayWIB()
	rows, err := h.Milestones.List(uid.(int64), repository.MilestoneFilter{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, m := range rows {
		if m.Status == "completed" || m.Status == "cancelled" {
			continue
		}
		if d := service.DaysUntil(m.TargetDate, today); d != nil && *d >= 0 && *d <= within {
			out = append(out, milestoneOut(m))
		}
	}
	c.JSON(http.StatusOK, out)
}
