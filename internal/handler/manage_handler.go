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

type SavingHandler struct {
	Savings *repository.SavingRepository
}

type savingIn struct {
	Name                string  `json:"name" binding:"required,max=200"`
	TargetAmount        float64 `json:"target_amount" binding:"required,gt=0"`
	CurrentAmount       float64 `json:"current_amount"`
	MonthlyContribution float64 `json:"monthly_contribution"`
	TargetDate          string  `json:"target_date"`
	Status              string  `json:"status"`
}

func savingOut(s *repository.SavingRow) gin.H {
	out := gin.H{"id": s.ID, "name": s.Name, "target_amount": s.TargetAmount,
		"current_amount": s.CurrentAmount, "monthly_contribution": s.MonthlyContribution,
		"progress":   service.GoalProgress(s.TargetAmount, s.CurrentAmount),
		"eta_months": service.EstimateMonths(s.TargetAmount, s.CurrentAmount, s.MonthlyContribution),
		"status":     s.Status}
	if s.TargetDate.Valid {
		out["target_date"] = s.TargetDate.String
	}
	return out
}

func toSavingRow(userID int64, in savingIn) (*repository.SavingRow, error) {
	st := in.Status
	if st == "" {
		st = "active"
	}
	if st != "active" && st != "paused" && st != "completed" && st != "cancelled" {
		return nil, errors.New("invalid")
	}
	s := &repository.SavingRow{UserID: userID, Name: in.Name, TargetAmount: in.TargetAmount,
		CurrentAmount: in.CurrentAmount, MonthlyContribution: in.MonthlyContribution, Status: st}
	if in.TargetDate != "" {
		s.TargetDate = sql.NullString{String: in.TargetDate, Valid: true}
	}
	return s, nil
}

func (h *SavingHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in savingIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	s, err := toSavingRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	id, err := h.Savings.Create(s)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	s.ID = id
	c.JSON(http.StatusCreated, savingOut(s))
}

func (h *SavingHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Savings.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, s := range rows {
		out = append(out, savingOut(s))
	}
	c.JSON(http.StatusOK, out)
}

func (h *SavingHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Savings.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in savingIn
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	s, err := toSavingRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	s.ID = cur.ID
	s.UserID = cur.UserID
	if err := h.Savings.Update(s); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, savingOut(s))
}

func (h *SavingHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Savings.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

type AssetHandler struct {
	Assets *repository.AssetRepository
}

func (h *AssetHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in struct {
		Name          string  `json:"name" binding:"required,max=200"`
		Category      string  `json:"category"`
		PurchasePrice float64 `json:"purchase_price"`
		CurrentValue  float64 `json:"current_value"`
		Condition     string  `json:"condition"`
		Location      string  `json:"location"`
		WarrantyEnd   string  `json:"warranty_end"`
		Notes         string  `json:"notes"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	cond := in.Condition
	if cond == "" {
		cond = "good"
	}
	if cond != "new" && cond != "good" && cond != "worn" && cond != "broken" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	a := &repository.AssetRow{UserID: uid.(int64), Name: in.Name, Category: in.Category,
		PurchasePrice: in.PurchasePrice, CurrentValue: in.CurrentValue,
		Condition: cond, Location: in.Location, Notes: in.Notes}
	if in.WarrantyEnd != "" {
		a.WarrantyEnd = sql.NullString{String: in.WarrantyEnd, Valid: true}
	}
	id, err := h.Assets.Create(a)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	a.ID = id
	c.JSON(http.StatusCreated, assetOut(a))
}

func assetOut(a *repository.AssetRow) gin.H {
	today := service.TodayWIB()
	out := gin.H{"id": a.ID, "name": a.Name, "category": a.Category,
		"purchase_price": a.PurchasePrice, "current_value": a.CurrentValue,
		"condition": a.Condition, "location": a.Location, "notes": a.Notes}
	if a.WarrantyEnd.Valid {
		out["warranty_end"] = a.WarrantyEnd.String
		out["warranty_expiring"] = service.DocExpiring(a.WarrantyEnd.String, today, 90)
	}
	return out
}

func (h *AssetHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Assets.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, a := range rows {
		out = append(out, assetOut(a))
	}
	c.JSON(http.StatusOK, out)
}

func (h *AssetHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Assets.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

type WishlistHandler struct {
	Wishlist *repository.WishlistRepository
}

func (h *WishlistHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in struct {
		Item       string  `json:"item" binding:"required,max=200"`
		Category   string  `json:"category"`
		Priority   string  `json:"priority"`
		EstPrice   float64 `json:"est_price"`
		Saved      float64 `json:"saved"`
		TargetDate string  `json:"target_date"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	pr := in.Priority
	if pr == "" {
		pr = "medium"
	}
	if pr != "low" && pr != "medium" && pr != "high" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	w := &repository.WishlistRow{UserID: uid.(int64), Item: in.Item, Category: in.Category,
		Priority: pr, EstPrice: in.EstPrice, Saved: in.Saved, Status: "planned"}
	if in.TargetDate != "" {
		w.TargetDate = sql.NullString{String: in.TargetDate, Valid: true}
	}
	id, err := h.Wishlist.Create(w)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	w.ID = id
	c.JSON(http.StatusCreated, wishlistOut(w))
}

func wishlistOut(w *repository.WishlistRow) gin.H {
	out := gin.H{"id": w.ID, "item": w.Item, "category": w.Category, "priority": w.Priority,
		"est_price": w.EstPrice, "saved": w.Saved,
		"progress": service.GoalProgress(w.EstPrice, w.Saved), "status": w.Status}
	if w.TargetDate.Valid {
		out["target_date"] = w.TargetDate.String
	}
	return out
}

func (h *WishlistHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Wishlist.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, w := range rows {
		out = append(out, wishlistOut(w))
	}
	c.JSON(http.StatusOK, out)
}

func (h *WishlistHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Wishlist.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

type DocumentHandler struct {
	Docs *repository.DocumentRepository
}

func (h *DocumentHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in struct {
		Item         string `json:"item" binding:"required,max=200"`
		Category     string `json:"category"`
		ExpiryDate   string `json:"expiry_date" binding:"required"`
		ReminderDays int    `json:"reminder_days"`
		Notes        string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	rd := in.ReminderDays
	if rd <= 0 {
		rd = 30
	}
	d := &repository.DocumentRow{UserID: uid.(int64), Item: in.Item, Category: in.Category,
		ExpiryDate: in.ExpiryDate, ReminderDays: rd, Notes: in.Notes}
	id, err := h.Docs.Create(d)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	d.ID = id
	c.JSON(http.StatusCreated, documentOut(d))
}

func documentOut(d *repository.DocumentRow) gin.H {
	today := service.TodayWIB()
	return gin.H{"id": d.ID, "item": d.Item, "category": d.Category, "expiry_date": d.ExpiryDate,
		"reminder_days": d.ReminderDays, "days_until": service.DaysUntil(d.ExpiryDate, today),
		"expiring": service.DocExpiring(d.ExpiryDate, today, d.ReminderDays), "notes": d.Notes}
}

func (h *DocumentHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Docs.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, d := range rows {
		out = append(out, documentOut(d))
	}
	c.JSON(http.StatusOK, out)
}

func (h *DocumentHandler) Upcoming(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Docs.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	today := service.TodayWIB()
	out := []gin.H{}
	for _, d := range rows {
		if service.DocExpiring(d.ExpiryDate, today, d.ReminderDays) {
			out = append(out, documentOut(d))
		}
	}
	c.JSON(http.StatusOK, out)
}

func (h *DocumentHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Docs.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

type ContactHandler struct {
	Contacts *repository.ContactRepository
}

func contactOut(x *repository.ContactRow) gin.H {
	today := service.TodayWIB()
	last := ""
	if x.LastContact.Valid {
		last = x.LastContact.String
	}
	out := gin.H{"id": x.ID, "name": x.Name, "relation": x.Relation,
		"last_contact": last, "contact_method": x.ContactMethod, "followup_days": x.FollowupDays,
		"followup_due": service.FollowupDue(last, today, x.FollowupDays), "notes": x.Notes}
	if x.Birthday.Valid {
		out["birthday"] = x.Birthday.String
	}
	return out
}

func (h *ContactHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in struct {
		Name          string `json:"name" binding:"required,max=200"`
		Relation      string `json:"relation"`
		Birthday      string `json:"birthday"`
		ContactMethod string `json:"contact_method"`
		FollowupDays  int    `json:"followup_days"`
		Notes         string `json:"notes"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		badRequest(c)
		return
	}
	fd := in.FollowupDays
	if fd <= 0 {
		fd = 30
	}
	x := &repository.ContactRow{UserID: uid.(int64), Name: in.Name, Relation: in.Relation,
		ContactMethod: in.ContactMethod, FollowupDays: fd, Notes: in.Notes}
	if in.Birthday != "" {
		x.Birthday = sql.NullString{String: in.Birthday, Valid: true}
	}
	id, err := h.Contacts.Create(x)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	x.ID = id
	c.JSON(http.StatusCreated, contactOut(x))
}

func (h *ContactHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Contacts.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, x := range rows {
		out = append(out, contactOut(x))
	}
	c.JSON(http.StatusOK, out)
}

func (h *ContactHandler) Touch(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if _, err := h.Contacts.Get(uid.(int64), id); errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	if err := h.Contacts.Touch(uid.(int64), id, service.TodayWIB()); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	x, _ := h.Contacts.Get(uid.(int64), id)
	c.JSON(http.StatusOK, contactOut(x))
}

func (h *ContactHandler) Followups(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Contacts.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, x := range rows {
		o := contactOut(x)
		if due, _ := o["followup_due"].(bool); due {
			out = append(out, o)
		}
	}
	c.JSON(http.StatusOK, out)
}

func (h *ContactHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if err := h.Contacts.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}
