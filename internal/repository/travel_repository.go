package repository

import "database/sql"

type TripRow struct {
	ID          int64
	UserID      int64
	Name        string
	Destination string
	StartDate   string
	EndDate     string
	Budget      float64
	Status      string
	Notes       string
}

type TripRepository struct{ DB *sql.DB }

const tripCols = `id,user_id,name,destination,start_date,end_date,budget,status,notes`

func (r *TripRepository) Create(t *TripRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO trips(user_id,name,destination,start_date,end_date,budget,status,notes)
		VALUES(?,?,?,?,?,?,?,?)`, t.UserID, t.Name, t.Destination, t.StartDate, t.EndDate, t.Budget, t.Status, t.Notes)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *TripRepository) Get(userID, id int64) (*TripRow, error) {
	var t TripRow
	err := r.DB.QueryRow(`SELECT `+tripCols+` FROM trips WHERE user_id=? AND id=?`, userID, id).
		Scan(&t.ID, &t.UserID, &t.Name, &t.Destination, &t.StartDate, &t.EndDate, &t.Budget, &t.Status, &t.Notes)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *TripRepository) List(userID int64) ([]*TripRow, error) {
	rows, err := r.DB.Query(`SELECT `+tripCols+` FROM trips WHERE user_id=? ORDER BY start_date`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*TripRow
	for rows.Next() {
		var t TripRow
		if err := rows.Scan(&t.ID, &t.UserID, &t.Name, &t.Destination, &t.StartDate, &t.EndDate, &t.Budget, &t.Status, &t.Notes); err != nil {
			return nil, err
		}
		out = append(out, &t)
	}
	return out, rows.Err()
}

func (r *TripRepository) Update(t *TripRow) error {
	_, err := r.DB.Exec(`UPDATE trips SET name=?,destination=?,start_date=?,end_date=?,budget=?,status=?,notes=?
		WHERE user_id=? AND id=?`, t.Name, t.Destination, t.StartDate, t.EndDate, t.Budget, t.Status, t.Notes, t.UserID, t.ID)
	return err
}

func (r *TripRepository) Delete(userID, id int64) error {
	_, err := r.DB.Exec(`DELETE FROM trips WHERE user_id=? AND id=?`, userID, id)
	return err
}

type ItineraryRow struct {
	ID       int64
	TripID   int64
	Date     string
	Time     string
	Activity string
	Location string
	Cost     float64
	Booked   int
}

type ItineraryRepository struct{ DB *sql.DB }

const itinCols = `id,trip_id,date,time,activity,location,cost,booked`

func (r *ItineraryRepository) Create(x *ItineraryRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO itinerary_items(trip_id,date,time,activity,location,cost,booked)
		VALUES(?,?,?,?,?,?,?)`, x.TripID, x.Date, x.Time, x.Activity, x.Location, x.Cost, x.Booked)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *ItineraryRepository) List(tripID int64) ([]*ItineraryRow, error) {
	rows, err := r.DB.Query(`SELECT `+itinCols+` FROM itinerary_items WHERE trip_id=? ORDER BY date,time`, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*ItineraryRow
	for rows.Next() {
		var x ItineraryRow
		if err := rows.Scan(&x.ID, &x.TripID, &x.Date, &x.Time, &x.Activity, &x.Location, &x.Cost, &x.Booked); err != nil {
			return nil, err
		}
		out = append(out, &x)
	}
	return out, rows.Err()
}

func (r *ItineraryRepository) SetBooked(id int64, booked bool) error {
	v := 0
	if booked {
		v = 1
	}
	_, err := r.DB.Exec(`UPDATE itinerary_items SET booked=? WHERE id=?`, v, id)
	return err
}

func (r *ItineraryRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM itinerary_items WHERE id=?`, id)
	return err
}

func (r *ItineraryRepository) TripCost(tripID int64) (float64, error) {
	var s sql.NullFloat64
	err := r.DB.QueryRow(`SELECT SUM(cost) FROM itinerary_items WHERE trip_id=?`, tripID).Scan(&s)
	if err != nil || !s.Valid {
		return 0, err
	}
	return s.Float64, nil
}

type PackingRow struct {
	ID       int64
	TripID   int64
	Category string
	Item     string
	Qty      int
	Packed   int
}

type PackingRepository struct{ DB *sql.DB }

const packCols = `id,trip_id,category,item,qty,packed`

func (r *PackingRepository) Create(x *PackingRow) (int64, error) {
	res, err := r.DB.Exec(`INSERT INTO packing_items(trip_id,category,item,qty,packed)
		VALUES(?,?,?,?,?)`, x.TripID, x.Category, x.Item, x.Qty, x.Packed)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func (r *PackingRepository) List(tripID int64) ([]*PackingRow, error) {
	rows, err := r.DB.Query(`SELECT `+packCols+` FROM packing_items WHERE trip_id=? ORDER BY id`, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() //nolint:errcheck
	var out []*PackingRow
	for rows.Next() {
		var x PackingRow
		if err := rows.Scan(&x.ID, &x.TripID, &x.Category, &x.Item, &x.Qty, &x.Packed); err != nil {
			return nil, err
		}
		out = append(out, &x)
	}
	return out, rows.Err()
}

func (r *PackingRepository) SetPacked(id int64, packed bool) error {
	v := 0
	if packed {
		v = 1
	}
	_, err := r.DB.Exec(`UPDATE packing_items SET packed=? WHERE id=?`, v, id)
	return err
}

func (r *PackingRepository) Delete(id int64) error {
	_, err := r.DB.Exec(`DELETE FROM packing_items WHERE id=?`, id)
	return err
}
