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

type TransactionHandler struct {
	Trx *repository.TransactionRepository
}

type trxIn struct {
	Date        string  `json:"date" binding:"required"`
	Type        string  `json:"type" binding:"required"`
	Category    string  `json:"category" binding:"required,max=100"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount" binding:"required,gt=0"`
	Account     string  `json:"account"`
	Recurring   *bool   `json:"recurring"`
}

func toTrxRow(userID int64, in trxIn) (*repository.TransactionRow, error) {
	if !service.ValidTrxType(in.Type) {
		return nil, errors.New("invalid")
	}
	if _, err := time.Parse("2006-01-02", in.Date); err != nil {
		return nil, errors.New("invalid")
	}
	t := &repository.TransactionRow{UserID: userID, Date: in.Date, Type: in.Type,
		Category: in.Category, Description: in.Description, Amount: in.Amount, Account: in.Account}
	if in.Recurring != nil && *in.Recurring {
		t.Recurring = 1
	}
	return t, nil
}

func trxOut(t *repository.TransactionRow) gin.H {
	return gin.H{"id": t.ID, "date": t.Date, "type": t.Type, "category": t.Category,
		"description": t.Description, "amount": t.Amount, "account": t.Account, "recurring": t.Recurring == 1}
}

func (h *TransactionHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in trxIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	t, err := toTrxRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	id, err := h.Trx.Create(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	t.ID = id
	c.JSON(http.StatusCreated, trxOut(t))
}

func (h *TransactionHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	var f repository.TrxFilter
	f.Type = c.Query("type")
	f.Category = c.Query("category")
	if y := c.Query("year"); y != "" {
		f.Year, _ = strconv.Atoi(y)
	}
	if m := c.Query("month"); m != "" {
		f.Month, _ = strconv.Atoi(m)
	}
	rows, err := h.Trx.List(uid.(int64), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, t := range rows {
		out = append(out, trxOut(t))
	}
	c.JSON(http.StatusOK, out)
}

func (h *TransactionHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	t, err := h.Trx.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, trxOut(t))
}

func (h *TransactionHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Trx.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in trxIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	t, err := toTrxRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	t.ID = cur.ID
	t.UserID = cur.UserID
	if err := h.Trx.Update(t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, trxOut(t))
}

func (h *TransactionHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Trx.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TransactionHandler) Summary(c *gin.Context) {
	uid, _ := c.Get("userID")
	y, _ := strconv.Atoi(c.DefaultQuery("year", service.TodayWIB()[:4]))
	m, _ := strconv.Atoi(c.Query("month"))
	if m < 1 || m > 12 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_month"})
		return
	}
	income, expense, err := h.Trx.MonthSummary(uid.(int64), y, m)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"year": y, "month": m,
		"income": income, "expense": expense, "net": income - expense})
}
