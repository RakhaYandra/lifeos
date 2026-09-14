package service

import "testing"

func TestOverdue(t *testing.T) {
	if !IsOverdue("in_progress", "2026-01-01", "2026-01-02") {
		t.Fatal("mau overdue")
	}
	if IsOverdue("completed", "2026-01-01", "2026-01-02") {
		t.Fatal("completed jangan overdue")
	}
	if IsOverdue("cancelled", "2026-01-01", "2026-01-02") {
		t.Fatal("cancelled jangan overdue")
	}
	if IsOverdue("inbox", "", "2026-01-02") {
		t.Fatal("tanpa due jangan overdue")
	}
	if IsOverdue("in_progress", "2026-01-02", "2026-01-02") {
		t.Fatal("due hari ini bukan overdue")
	}
}

func TestDaysRemaining(t *testing.T) {
	n := DaysRemaining("2026-01-05", "2026-01-02")
	if n == nil || *n != 3 {
		t.Fatalf("mau 3, dapat %v", n)
	}
	n = DaysRemaining("2026-01-01", "2026-01-02")
	if n == nil || *n != -1 {
		t.Fatalf("mau -1, dapat %v", n)
	}
	if DaysRemaining("", "2026-01-02") != nil {
		t.Fatal("tanpa due mau nil")
	}
	if DaysRemaining("salah", "2026-01-02") != nil {
		t.Fatal("tanggal rusak mau nil")
	}
}

func TestProjectProgress(t *testing.T) {
	if ProjectProgress(0, 0) != 0 {
		t.Fatal("total 0 mau 0")
	}
	if ProjectProgress(4, 1) != 25 {
		t.Fatal("mau 25")
	}
	if ProjectProgress(2, 5) != 100 {
		t.Fatal("mau cap 100")
	}
}
