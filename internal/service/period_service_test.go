package service

import "testing"

func TestMonthBounds(t *testing.T) {
	f, to, ok := MonthBounds("2026-09")
	if !ok || f != "2026-09-01" || to != "2026-09-30" {
		t.Fatalf("mau 01-30, dapat %s-%s", f, to)
	}
	f, to, ok = MonthBounds("2026-02")
	if !ok || f != "2026-02-01" || to != "2026-02-28" {
		t.Fatalf("mau feb, dapat %s-%s", f, to)
	}
	if _, _, ok := MonthBounds("2026-13"); ok {
		t.Fatal("bulan 13 tolak")
	}
	if _, _, ok := MonthBounds("2026"); ok {
		t.Fatal("tanpa bulan tolak")
	}
}

func TestYearBounds(t *testing.T) {
	f, to, ok := YearBounds("2026")
	if !ok || f != "2026-01-01" || to != "2026-12-31" {
		t.Fatalf("mau setahun, dapat %s-%s", f, to)
	}
	if _, _, ok := YearBounds("26"); ok {
		t.Fatal("format pendek tolak")
	}
}

func TestCountInRange(t *testing.T) {
	if CountInRange([]string{"2026-09-01", "2026-09-15", "2026-10-01"}, "2026-09-01", "2026-09-30") != 2 {
		t.Fatal("mau 2")
	}
}
