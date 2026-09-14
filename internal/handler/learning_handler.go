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

type LearningHandler struct {
	Learning *repository.LearningRepository
	Reading  *repository.ReadingRepository
}

type learnIn struct {
	Topic        string  `json:"topic" binding:"required,max=300"`
	Type         string  `json:"type"`
	Provider     string  `json:"provider"`
	RelatedSkill string  `json:"related_skill"`
	Progress     int     `json:"progress"`
	Hours        float64 `json:"hours"`
	Status       string  `json:"status"`
	Notes        string  `json:"notes"`
}

func toLearnRow(userID int64, in learnIn) (*repository.LearningRow, error) {
	ty := in.Type
	if ty == "" {
		ty = "course"
	}
	st := in.Status
	if st == "" {
		st = "active"
	}
	if !service.ValidLearnType(ty) || !service.ValidLearnStatus(st) {
		return nil, errors.New("invalid")
	}
	if in.Progress < 0 || in.Progress > 100 {
		return nil, errors.New("invalid")
	}
	return &repository.LearningRow{UserID: userID, Topic: in.Topic, Type: ty, Provider: in.Provider,
		RelatedSkill: in.RelatedSkill, Progress: in.Progress, Hours: in.Hours, Status: st, Notes: in.Notes}, nil
}

func learnOut(l *repository.LearningRow) gin.H {
	return gin.H{"id": l.ID, "topic": l.Topic, "type": l.Type, "provider": l.Provider,
		"related_skill": l.RelatedSkill, "progress": l.Progress, "hours": l.Hours, "status": l.Status, "notes": l.Notes}
}

func (h *LearningHandler) CreateLearn(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in learnIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	l, err := toLearnRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	id, err := h.Learning.Create(l)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	l.ID = id
	c.JSON(http.StatusCreated, learnOut(l))
}

func (h *LearningHandler) ListLearn(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Learning.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, l := range rows {
		out = append(out, learnOut(l))
	}
	c.JSON(http.StatusOK, out)
}

func (h *LearningHandler) UpdateLearn(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Learning.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in learnIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	l, err := toLearnRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	l.ID = cur.ID
	l.UserID = cur.UserID
	if err := h.Learning.Update(l); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, learnOut(l))
}

func (h *LearningHandler) DeleteLearn(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Learning.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

type readIn struct {
	Title     string `json:"title" binding:"required,max=300"`
	Type      string `json:"type"`
	Author    string `json:"author"`
	Status    string `json:"status"`
	Rating    *int   `json:"rating"`
	Takeaways string `json:"takeaways"`
}

func toReadRow(userID int64, in readIn) (*repository.ReadingRow, error) {
	ty := in.Type
	if ty == "" {
		ty = "book"
	}
	st := in.Status
	if st == "" {
		st = "reading"
	}
	if !service.ValidReadType(ty) || !service.ValidReadStatus(st) {
		return nil, errors.New("invalid")
	}
	x := &repository.ReadingRow{UserID: userID, Title: in.Title, Type: ty, Author: in.Author, Status: st, Takeaways: in.Takeaways}
	if in.Rating != nil {
		if *in.Rating < 1 || *in.Rating > 5 {
			return nil, errors.New("invalid")
		}
		x.Rating = sql.NullInt64{Int64: int64(*in.Rating), Valid: true}
	}
	return x, nil
}

func readOut(x *repository.ReadingRow) gin.H {
	out := gin.H{"id": x.ID, "title": x.Title, "type": x.Type, "author": x.Author,
		"status": x.Status, "takeaways": x.Takeaways}
	if x.Rating.Valid {
		out["rating"] = x.Rating.Int64
	}
	return out
}

func (h *LearningHandler) CreateRead(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in readIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	x, err := toReadRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	id, err := h.Reading.Create(x)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	x.ID = id
	c.JSON(http.StatusCreated, readOut(x))
}

func (h *LearningHandler) ListRead(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Reading.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, x := range rows {
		out = append(out, readOut(x))
	}
	c.JSON(http.StatusOK, out)
}

func (h *LearningHandler) UpdateRead(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Reading.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in readIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	x, err := toReadRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	x.ID = cur.ID
	x.UserID = cur.UserID
	if err := h.Reading.Update(x); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, readOut(x))
}

func (h *LearningHandler) DeleteRead(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Reading.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}
