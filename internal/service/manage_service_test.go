package service

import "testing"

func TestEstimateMonths(t *testing.T) {
	if EstimateMonths(20000000, 8000000, 800000) != 15 {
		t.Fatal("mau 15")
	}
	if EstimateMonths(100, 100, 10) != -1 {
		t.Fatal("lunas mau -1")
	}
	if EstimateMonths(100, 50, 0) != -1 {
		t.Fatal("tanpa kontribusi mau -1")
	}
	if EstimateMonths(100, 0, 30) != 4 {
		t.Fatal("mau ceil 4")
	}
}

func TestFollowupDue(t *testing.T) {
	if !FollowupDue("2026-08-01", "2026-09-14", 30) {
		t.Fatal("44 hari mau due")
	}
	if FollowupDue("2026-09-10", "2026-09-14", 30) {
		t.Fatal("4 hari jangan due")
	}
	if !FollowupDue("", "2026-09-14", 30) {
		t.Fatal("belum pernah mau due")
	}
	if FollowupDue("2026-08-01", "2026-09-14", 0) {
		t.Fatal("followup 0 jangan due")
	}
}

func TestDocExpiring(t *testing.T) {
	if !DocExpiring("2026-09-20", "2026-09-14", 30) {
		t.Fatal("6 hari mau expiring")
	}
	if DocExpiring("2027-02-15", "2026-09-14", 30) {
		t.Fatal("jauh jangan expiring")
	}
	if !DocExpiring("2026-09-01", "2026-09-14", 30) {
		t.Fatal("lewat mau expiring")
	}
}
