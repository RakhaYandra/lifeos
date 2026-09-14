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

type ProjectHandler struct {
	Projects *repository.ProjectRepository
}

type projectIn struct {
	Name       string `json:"name" binding:"required,max=200"`
	LifeAreaID *int64 `json:"life_area_id"`
	GoalID     *int64 `json:"goal_id"`
	Status     string `json:"status"`
	StartDate  string `json:"start_date"`
	TargetDate string `json:"target_date"`
	Notes      string `json:"notes"`
}

func toProjectRow(userID int64, in projectIn) (*repository.ProjectRow, error) {
	st := in.Status
	if st == "" {
		st = "active"
	}
	if !service.ValidProjectStatus(st) {
		return nil, errors.New("status_invalid")
	}
	p := &repository.ProjectRow{UserID: userID, Name: in.Name, Status: st, Notes: in.Notes}
	if in.LifeAreaID != nil {
		p.LifeAreaID = sql.NullInt64{Int64: *in.LifeAreaID, Valid: true}
	}
	if in.GoalID != nil {
		p.GoalID = sql.NullInt64{Int64: *in.GoalID, Valid: true}
	}
	if in.StartDate != "" {
		p.StartDate = sql.NullString{String: in.StartDate, Valid: true}
	}
	if in.TargetDate != "" {
		p.TargetDate = sql.NullString{String: in.TargetDate, Valid: true}
	}
	return p, nil
}

func projectOut(p *repository.ProjectRow, total, done int) gin.H {
	out := gin.H{"id": p.ID, "name": p.Name, "status": p.Status,
		"progress":    service.ProjectProgress(total, done),
		"total_tasks": total, "done_tasks": done, "notes": p.Notes}
	if p.LifeAreaID.Valid {
		out["life_area_id"] = p.LifeAreaID.Int64
	}
	if p.GoalID.Valid {
		out["goal_id"] = p.GoalID.Int64
	}
	if p.StartDate.Valid {
		out["start_date"] = p.StartDate.String
	}
	if p.TargetDate.Valid {
		out["target_date"] = p.TargetDate.String
	}
	if p.CompletedAt.Valid {
		out["completed_at"] = p.CompletedAt.String
	}
	return out
}

func (h *ProjectHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in projectIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	p, err := toProjectRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status_invalid"})
		return
	}
	id, err := h.Projects.Create(p)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	p.ID = id
	c.JSON(http.StatusCreated, projectOut(p, 0, 0))
}

func (h *ProjectHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Projects.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, p := range rows {
		total, done, _ := h.Projects.TaskCounts(p.ID)
		out = append(out, projectOut(p, total, done))
	}
	c.JSON(http.StatusOK, out)
}

func (h *ProjectHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	p, err := h.Projects.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	total, done, _ := h.Projects.TaskCounts(p.ID)
	c.JSON(http.StatusOK, projectOut(p, total, done))
}

func (h *ProjectHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Projects.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in projectIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	p, err := toProjectRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "status_invalid"})
		return
	}
	p.ID = cur.ID
	p.UserID = cur.UserID
	p.CompletedAt = cur.CompletedAt
	if p.Status == "completed" && !cur.CompletedAt.Valid {
		p.CompletedAt = sql.NullString{String: service.TodayWIB(), Valid: true}
	}
	if err := h.Projects.Update(p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	total, done, _ := h.Projects.TaskCounts(p.ID)
	c.JSON(http.StatusOK, projectOut(p, total, done))
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Projects.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}
