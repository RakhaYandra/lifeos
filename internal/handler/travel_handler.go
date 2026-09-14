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

type TravelHandler struct {
	Trips *repository.TripRepository
	Itin  *repository.ItineraryRepository
	Pack  *repository.PackingRepository
}

type tripIn struct {
	Name        string  `json:"name" binding:"required,max=300"`
	Destination string  `json:"destination"`
	StartDate   string  `json:"start_date" binding:"required"`
	EndDate     string  `json:"end_date" binding:"required"`
	Budget      float64 `json:"budget"`
	Status      string  `json:"status"`
	Notes       string  `json:"notes"`
}

func toTripRow(userID int64, in tripIn) (*repository.TripRow, error) {
	st := in.Status
	if st == "" {
		st = "planning"
	}
	if !service.ValidTripStatus(st) {
		return nil, errors.New("invalid")
	}
	s, err1 := time.Parse("2006-01-02", in.StartDate)
	e, err2 := time.Parse("2006-01-02", in.EndDate)
	if err1 != nil || err2 != nil || e.Before(s) {
		return nil, errors.New("invalid")
	}
	return &repository.TripRow{UserID: userID, Name: in.Name, Destination: in.Destination,
		StartDate: in.StartDate, EndDate: in.EndDate, Budget: in.Budget, Status: st, Notes: in.Notes}, nil
}

func (h *TravelHandler) tripOut(t *repository.TripRow) gin.H {
	cost, _ := h.Itin.TripCost(t.ID)
	return gin.H{"id": t.ID, "name": t.Name, "destination": t.Destination,
		"start_date": t.StartDate, "end_date": t.EndDate, "budget": t.Budget,
		"actual_cost": cost, "status": t.Status, "notes": t.Notes}
}

func (h *TravelHandler) Create(c *gin.Context) {
	uid, _ := c.Get("userID")
	var in tripIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := toTripRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	id, err := h.Trips.Create(t)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	t.ID = id
	c.JSON(http.StatusCreated, h.tripOut(t))
}

func (h *TravelHandler) List(c *gin.Context) {
	uid, _ := c.Get("userID")
	rows, err := h.Trips.List(uid.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	out := []gin.H{}
	for _, t := range rows {
		out = append(out, h.tripOut(t))
	}
	c.JSON(http.StatusOK, out)
}

func (h *TravelHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	t, err := h.Trips.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	itin, _ := h.Itin.List(t.ID)
	pack, _ := h.Pack.List(t.ID)
	io := []gin.H{}
	for _, x := range itin {
		io = append(io, gin.H{"id": x.ID, "date": x.Date, "time": x.Time, "activity": x.Activity,
			"location": x.Location, "cost": x.Cost, "booked": x.Booked == 1})
	}
	po := []gin.H{}
	for _, x := range pack {
		po = append(po, gin.H{"id": x.ID, "category": x.Category, "item": x.Item, "qty": x.Qty, "packed": x.Packed == 1})
	}
	out := h.tripOut(t)
	out["itinerary"] = io
	out["packing"] = po
	c.JSON(http.StatusOK, out)
}

func (h *TravelHandler) Update(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	cur, err := h.Trips.Get(uid.(int64), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	var in tripIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	t, err := toTripRow(uid.(int64), in)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid"})
		return
	}
	t.ID = cur.ID
	t.UserID = cur.UserID
	if err := h.Trips.Update(t); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusOK, h.tripOut(t))
}

func (h *TravelHandler) Delete(c *gin.Context) {
	uid, _ := c.Get("userID")
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if _, err := h.Trips.Get(uid.(int64), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if err := h.Trips.Delete(uid.(int64), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TravelHandler) ownedTrip(uid, id int64) (*repository.TripRow, bool) {
	t, err := h.Trips.Get(uid, id)
	if err != nil {
		return nil, false
	}
	return t, true
}

type itinIn struct {
	Date     string  `json:"date" binding:"required"`
	Time     string  `json:"time"`
	Activity string  `json:"activity" binding:"required,max=300"`
	Location string  `json:"location"`
	Cost     float64 `json:"cost"`
}

func (h *TravelHandler) AddItin(c *gin.Context) {
	uid, _ := c.Get("userID")
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if _, ok := h.ownedTrip(uid.(int64), tid); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	var in itinIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if _, err := time.Parse("2006-01-02", in.Date); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_date"})
		return
	}
	id, err := h.Itin.Create(&repository.ItineraryRow{TripID: tid, Date: in.Date, Time: in.Time,
		Activity: in.Activity, Location: in.Location, Cost: in.Cost})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *TravelHandler) ToggleItin(c *gin.Context) {
	uid, _ := c.Get("userID")
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	iid, _ := strconv.ParseInt(c.Param("iid"), 10, 64)
	if _, ok := h.ownedTrip(uid.(int64), tid); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	var in struct {
		Booked bool `json:"booked"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Itin.SetBooked(iid, in.Booked); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TravelHandler) DeleteItin(c *gin.Context) {
	uid, _ := c.Get("userID")
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	iid, _ := strconv.ParseInt(c.Param("iid"), 10, 64)
	if _, ok := h.ownedTrip(uid.(int64), tid); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if err := h.Itin.Delete(iid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

type packIn struct {
	Category string `json:"category"`
	Item     string `json:"item" binding:"required,max=200"`
	Qty      int    `json:"qty"`
}

func (h *TravelHandler) AddPack(c *gin.Context) {
	uid, _ := c.Get("userID")
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	if _, ok := h.ownedTrip(uid.(int64), tid); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	var in packIn
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	qty := in.Qty
	if qty <= 0 {
		qty = 1
	}
	id, err := h.Pack.Create(&repository.PackingRow{TripID: tid, Category: in.Category, Item: in.Item, Qty: qty})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id})
}

func (h *TravelHandler) TogglePack(c *gin.Context) {
	uid, _ := c.Get("userID")
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	pid, _ := strconv.ParseInt(c.Param("pid"), 10, 64)
	if _, ok := h.ownedTrip(uid.(int64), tid); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	var in struct {
		Packed bool `json:"packed"`
	}
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Pack.SetPacked(pid, in.Packed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *TravelHandler) DeletePack(c *gin.Context) {
	uid, _ := c.Get("userID")
	tid, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	pid, _ := strconv.ParseInt(c.Param("pid"), 10, 64)
	if _, ok := h.ownedTrip(uid.(int64), tid); !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "not_found"})
		return
	}
	if err := h.Pack.Delete(pid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal"})
		return
	}
	c.Status(http.StatusNoContent)
}
