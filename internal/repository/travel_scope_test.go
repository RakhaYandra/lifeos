package repository

import (
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

// Skema minimal untuk menguji scoping kepemilikan parent-child.
const testTravelSchema = `
CREATE TABLE itinerary_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  trip_id INTEGER NOT NULL,
  activity TEXT NOT NULL DEFAULT '',
  booked INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE packing_items (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  trip_id INTEGER NOT NULL,
  item TEXT NOT NULL DEFAULT '',
  packed INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE decision_options (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  decision_id INTEGER NOT NULL,
  name TEXT NOT NULL DEFAULT ''
);
CREATE TABLE decision_marks (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  option_id INTEGER NOT NULL,
  criterion TEXT NOT NULL,
  weight REAL NOT NULL DEFAULT 0,
  score REAL NOT NULL DEFAULT 0,
  UNIQUE(option_id, criterion)
);
`

func newTestDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	if _, err := db.Exec(testTravelSchema); err != nil {
		t.Fatalf("schema: %v", err)
	}
	t.Cleanup(func() { db.Close() }) //nolint:errcheck
	return db
}

// -IDOR: SetBooked/Delete hanya boleh menyentuh baris milik trip yang diberikan.
func TestItineraryScopedToTrip(t *testing.T) {
	db := newTestDB(t)
	r := &ItineraryRepository{DB: db}

	if _, err := db.Exec(`INSERT INTO itinerary_items (id, trip_id, activity) VALUES (1, 100, 'milik orang lain')`); err != nil {
		t.Fatal(err)
	}

	// Trip 999 bukan pemilik item 1.
	ok, err := r.SetBooked(999, 1, true)
	if err != nil {
		t.Fatalf("SetBooked: %v", err)
	}
	if ok {
		t.Fatal("SetBooked mengubah baris milik trip lain — IDOR belum tertutup")
	}

	var booked int
	if err := db.QueryRow(`SELECT booked FROM itinerary_items WHERE id=1`).Scan(&booked); err != nil {
		t.Fatal(err)
	}
	if booked != 0 {
		t.Fatalf("baris tetap berubah: booked=%d", booked)
	}

	// Pemilik sah tetap bisa mengubah.
	ok, err = r.SetBooked(100, 1, true)
	if err != nil || !ok {
		t.Fatalf("pemilik sah ditolak: ok=%v err=%v", ok, err)
	}
	if err := db.QueryRow(`SELECT booked FROM itinerary_items WHERE id=1`).Scan(&booked); err != nil {
		t.Fatal(err)
	}
	if booked != 1 {
		t.Fatalf("pemilik sah tidak bisa mengubah: booked=%d", booked)
	}
}

func TestItineraryDeleteScopedToTrip(t *testing.T) {
	db := newTestDB(t)
	r := &ItineraryRepository{DB: db}
	if _, err := db.Exec(`INSERT INTO itinerary_items (id, trip_id) VALUES (1, 100)`); err != nil {
		t.Fatal(err)
	}

	ok, err := r.Delete(999, 1)
	if err != nil || ok {
		t.Fatalf("Delete mengizinkan trip lain: ok=%v err=%v", ok, err)
	}
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM itinerary_items WHERE id=1`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("baris terhapus oleh trip yang tidak berhak")
	}

	if ok, err := r.Delete(100, 1); err != nil || !ok {
		t.Fatalf("pemilik sah tidak bisa hapus: ok=%v err=%v", ok, err)
	}
}

func TestPackingScopedToTrip(t *testing.T) {
	db := newTestDB(t)
	r := &PackingRepository{DB: db}
	if _, err := db.Exec(`INSERT INTO packing_items (id, trip_id) VALUES (1, 100)`); err != nil {
		t.Fatal(err)
	}

	if ok, err := r.SetPacked(999, 1, true); err != nil || ok {
		t.Fatalf("SetPacked mengizinkan trip lain: ok=%v err=%v", ok, err)
	}
	if ok, err := r.Delete(999, 1); err != nil || ok {
		t.Fatalf("Delete mengizinkan trip lain: ok=%v err=%v", ok, err)
	}

	if ok, err := r.SetPacked(100, 1, true); err != nil || !ok {
		t.Fatalf("pemilik sah ditolak: ok=%v err=%v", ok, err)
	}
}

func TestDecisionOptionScopedToDecision(t *testing.T) {
	db := newTestDB(t)
	r := &DecisionOptionRepository{DB: db}
	if _, err := db.Exec(`INSERT INTO decision_options (id, decision_id, name) VALUES (1, 100, 'milik orang lain')`); err != nil {
		t.Fatal(err)
	}

	ok, err := r.Delete(999, 1)
	if err != nil || ok {
		t.Fatalf("Delete mengizinkan decision lain: ok=%v err=%v", ok, err)
	}
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM decision_options WHERE id=1`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatal("opsi terhapus oleh decision yang tidak berhak")
	}

	if ok, err := r.Delete(100, 1); err != nil || !ok {
		t.Fatalf("pemilik sah tidak bisa hapus: ok=%v err=%v", ok, err)
	}
}

func TestDecisionMarkScopedToDecision(t *testing.T) {
	db := newTestDB(t)
	marks := &DecisionMarkRepository{DB: db}
	if _, err := db.Exec(`INSERT INTO decision_options (id, decision_id) VALUES (1, 100)`); err != nil {
		t.Fatal(err)
	}

	// Opsi 1 milik decision 100; decision 999 tidak berhak menandai.
	ok, err := marks.Upsert(999, 1, "harga", 1, 5)
	if err != nil {
		t.Fatalf("Upsert: %v", err)
	}
	if ok {
		t.Fatal("Upsert menandai opsi milik decision lain — IDOR belum tertutup")
	}
	var n int
	if err := db.QueryRow(`SELECT COUNT(*) FROM decision_marks`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatal("mark tetap tertulis untuk decision yang tidak berhak")
	}

	// Opsi yang tidak ada sama sekali juga ditolak.
	if ok, err := marks.Upsert(100, 4242, "harga", 1, 5); err != nil || ok {
		t.Fatalf("Upsert menerima opsi tak ada: ok=%v err=%v", ok, err)
	}

	// Pemilik sah bisa.
	if ok, err := marks.Upsert(100, 1, "harga", 1, 5); err != nil || !ok {
		t.Fatalf("pemilik sah ditolak: ok=%v err=%v", ok, err)
	}
}
