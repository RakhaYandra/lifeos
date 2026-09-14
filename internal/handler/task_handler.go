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

type TaskHandler struct {
	Tasks *repository.TaskRepository
}

type taskIn struct {
	ProjectID    *int64   `json:"project_id"`
	GoalID       *int64   `json:"goal_id"`
	Title        string   `json:"title" binding:"required,max=300"`
	Description  string   `json:"description"`
	LifeAreaID   *int64   `json:"life_area_id"`
	Status       string   `json:"status"`
	Priority     string   `json:"priority"`
	DueDate      string   `json:"due_date"`
	EffortEst    *float64 `json:"effort_est"`
	EffortActual *float64 `json:"effort_actual"`
}

func toTaskRow(userID int64, in taskIn) (*repository.TaskRow, error) {
	st := in.Status
	if st == "" {
		st = "inbox"
	}
	pr := in.Priority
	if pr == "" {
		pr = "medium"
	}
	if !service.ValidTaskStatus(st) || !service.ValidTaskPriority(pr) {
		return nil, errors.New("invalid")
	}
	t := &repository.TaskRow{UserID: userID, Title: in.Title, Description: in.Description, Status: st, Priority: pr}
	if in.ProjectID != nil {
		t.ProjectID = sql.NullInt64{Int64: *in.ProjectID, Valid: true}
	}
	if in.GoalID != nil {
		t.GoalID = sql.NullInt64{Int64: *in.GoalID, Valid: true}
	}
	if in.LifeAreaID != nil {
		t.LifeAreaID = sql.NullInt64{Int64: *in.LifeAreaID, Valid: true}
	}
	if in.DueDate != "" {
		if _, err := time.Parse("2006-01-02", in.DueDate); err != nil {
			return nil, errors.New("invalid")
		}
		t.DueDate = sql.NullString{String: in.DueDate, Valid: true}
	}
	if in.EffortEst != nil {
		t.EffortEst = sql.NullFloat64{Float64: *in.EffortEst, Valid: true}
	}
	if in.EffortActual != nil {
		t.EffortActual = sql.NullFloat64{Float64: *in.EffortActual, Valid: true}
	}
	return t, nil
}

func taskOut(t *repository.TaskRow, today string) gin.H {
	due := ""
	if t.DueDate.Valid {
		due = t.DueDate.String
	}
	out := gin.H{"id": t.ID, "title": t.Title, "description": t.Description,
		"status": t.Status, "priority": t.Priority,
		"overdue":        service.IsOverdue(t.Status, due, today),
		"days_remaining": service.DaysRemaining(due, today)}
	if t.ProjectID.Valid {
		out["project_id"] = t.ProjectID.Int64
	}
	if t.GoalID.Valid {
		out["goal_id"] = t.GoalID.Int64
	}
	if t.LifeAreaID.Valid {
		out["life_area_id"] = t.LifeAreaID.Int64
	}
	if due != "" {
		out["due_date"] = due
	}
	if t.CompletedAt.Valid {
		out["completed_at"] = t.CompletedAt.String
	}
	if t.EffortEst.Valid {
		out["effort_est"] = t.EffortEst.Float64
	}
	if t.EffortActual.Valid {
		out["effort_actual"] = t.EffortActual.Float64
	}
	return out
}

func (h *TaskHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in taskIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	t, err := toTaskRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	id, err := h.Tasks.Create(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	t.ID = id
	c.JSON(http.StatusCreated, taskOut(t, service.TodayWIB()))
}

func (h *TaskHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	var f repository.TaskFilter
	f.Status = c.Query("status")
	if pid := c.Query("project_id"); pid != "" {
		f.ProjectID, _ = strconv.ParseInt(pid, 10, 64)
	}
	rows, err := h.Tasks.List(uid.(int64), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	today := service.TodayWIB()
	out := []gin.H{}
	for _, t := range rows {
		out = append(out, taskOut(t, today))
	}
	c.JSON(http.StatusOK, out)
}

func (h *TaskHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	t, err := h.Tasks.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, taskOut(t, service.TodayWIB()))
}

func (h *TaskHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Tasks.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in taskIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	t, err := toTaskRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	t.ID = cur.ID
	t.UserID = cur.UserID
	t.CompletedAt = cur.CompletedAt
	if t.Status == "completed" && !cur.CompletedAt.Valid {
		t.CompletedAt = sql.NullString{String: service.TodayWIB(), Valid: true}
	}
	if t.Status != "completed" && cur.Status == "completed" {
		t.CompletedAt = sql.NullString{}
	}
	if err := h.Tasks.Update(t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, taskOut(t, service.TodayWIB()))
}

func (h *TaskHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Tasks.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TaskHandler) Today(c *gin.Context) {
	uid, _ := c.Get("userID")
	today := service.TodayWIB()
	rows, err := h.Tasks.List(uid.(int64), repository.TaskFilter{DueFrom: today, DueTo: today})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, t := range rows {
		out = append(out, taskOut(t, today))
	}
	c.JSON(http.StatusOK, out)
}

func (h *TaskHandler) Week(c *gin.Context) {
	uid, _ := c.Get("userID")
	start := c.Query("start")
	if start == "" {
		start = service.TodayWIB()
	}
	s, err := time.Parse("2006-01-02", start)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_start"})
		return
	}
	end := s.AddDate(0, 0, 6).Format("2006-01-02")
	rows, err := h.Tasks.List(uid.(int64), repository.TaskFilter{DueFrom: start, DueTo: end})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	today := service.TodayWIB()
	out := []gin.H{}
	for _, t := range rows {
		out = append(out, taskOut(t, today))
	}
	c.JSON(http.StatusOK, gin.H{"start": start, "end": end, "tasks": out})
}
