package service

import "testing"

func TestWeekRange(t *testing.T) {
	s, e := WeekRange("2026-09-14") // Senin
	if s != "2026-09-14" || e != "2026-09-20" {
		t.Fatalf("mau 14-20, dapat %s-%s", s, e)
	}
	s, e = WeekRange("2026-09-13") // Minggu → pekan sama
	if s != "2026-09-07" || e != "2026-09-13" {
		t.Fatalf("mau 07-13, dapat %s-%s", s, e)
	}
}

func TestNextOccurrence(t *testing.T) {
	if NextOccurrence("2026-09-10", "none", "2026-09-14") != "2026-09-10" {
		t.Fatal("none jangan geser")
	}
	if NextOccurrence("2026-09-10", "daily", "2026-09-14") != "2026-09-15" {
		t.Fatal("daily mau besok")
	}
	if NextOccurrence("2026-09-01", "weekly", "2026-09-14") != "2026-09-15" {
		t.Fatal("weekly mau 15")
	}
	if NextOccurrence("2026-09-20", "monthly", "2026-09-14") != "2026-09-20" {
		t.Fatal("masa depan jangan geser")
	}
	if NextOccurrence("2020-01-15", "yearly", "2026-09-14") != "2027-01-15" {
		t.Fatal("yearly mau 2027")
	}
}
