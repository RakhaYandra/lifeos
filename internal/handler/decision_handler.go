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

type DecisionHandler struct {
	Decisions *repository.DecisionRepository
	Options   *repository.DecisionOptionRepository
	Marks     *repository.DecisionMarkRepository
}

func (h *DecisionHandler) detail(uid, id int64) (gin.H, bool) {
	d, err := h.Decisions.Get(uid, id)
	if err != nil {
		return nil, false
	}
	opts, _ := h.Options.List(d.ID)
	marks, _ := h.Marks.ListByDecision(d.ID)
	byOpt := map[int64][]gin.H{}
	totals := map[int64]float64{}
	for _, m := range marks {
		byOpt[m.OptionID] = append(byOpt[m.OptionID],
			gin.H{"criterion": m.Criterion, "weight": m.Weight, "score": m.Score})
		totals[m.OptionID] += m.Weight * m.Score
	}
	ranks := service.RankOptions(totals)
	out := []gin.H{}
	for _, o := range opts {
		marks := byOpt[o.ID]
		if marks == nil {
			marks = []gin.H{}
		}
		out = append(out, gin.H{"id": o.ID, "name": o.Name, "marks": marks,
			"total": totals[o.ID], "rank": ranks[o.ID]})
	}
	return gin.H{"id": d.ID, "title": d.Title, "notes": d.Notes, "options": out}, true
}

func (h *DecisionHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in struct {
		Title string `json:"title" binding:"required,max=300"`
		Notes string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	id, err := h.Decisions.Create(&repository.DecisionRow{UserID: uid.(int64), Title: in.Title, Notes: in.Notes})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out, _ := h.detail(uid.(int64), id)
	c.JSON(http.StatusCreated, out)
}

func (h *DecisionHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Decisions.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, d := range rows {
		o, _ := h.detail(uid.(int64), d.ID)
		out = append(out, o)
	}
	c.JSON(http.StatusOK, out)
}

func (h *DecisionHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	out, ok := h.detail(uid.(int64), id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	c.JSON(http.StatusOK, out)
}

func (h *DecisionHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if _, err := h.Decisions.Get(uid.(int64), id); errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if err := h.Decisions.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DecisionHandler) AddOption(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if _, err := h.Decisions.Get(uid.(int64), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	var in struct {
		Name string `json:"name" binding:"required,max=200"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	oid, err := h.Options.Create(id, in.Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": oid, "name": in.Name})
}

func (h *DecisionHandler) DeleteOption(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	oid, _ := strconv.ParseInt(c.Param("oid"), 10, 64)
	if _, err := h.Decisions.Get(uid.(int64), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if err := h.Options.Delete(oid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *DecisionHandler) Mark(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	oid, _ := strconv.ParseInt(c.Param("oid"), 10, 64)
	if _, err := h.Decisions.Get(uid.(int64), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	var in struct {
		Criterion string  `json:"criterion" binding:"required,max=100"`
		Weight    float64 `json:"weight"`
		Score     float64 `json:"score"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	if in.Weight <= 0 {
		in.Weight = 1
	}
	if err := h.Marks.Upsert(oid, in.Criterion, in.Weight, in.Score); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out, _ := h.detail(uid.(int64), id)
	c.JSON(http.StatusOK, out)
}
