package service

import "testing"

func TestUtilization(t *testing.T) {
	if Utilization(50, 100) != 50 {
		t.Fatal("mau 50")
	}
	if Utilization(10, 0) != 0 {
		t.Fatal("budget 0 mau 0")
	}
}

func TestBudgetStatus(t *testing.T) {
	if BudgetStatus(50, 80) != "safe" {
		t.Fatal("mau safe")
	}
	if BudgetStatus(85, 80) != "warning" {
		t.Fatal("mau warning")
	}
	if BudgetStatus(100, 80) != "over" {
		t.Fatal("100 mau over")
	}
	if BudgetStatus(150, 80) != "over" {
		t.Fatal("mau over")
	}
}

func TestAnnualCost(t *testing.T) {
	if AnnualCost(50000, "monthly") != 600000 {
		t.Fatal("mau 600rb")
	}
	if AnnualCost(500000, "yearly") != 500000 {
		t.Fatal("mau 500rb")
	}
	if AnnualCost(10000, "weekly") != 520000 {
		t.Fatal("mau 520rb")
	}
}

func TestDaysUntil(t *testing.T) {
	n := DaysUntil("2026-09-20", "2026-09-14")
	if n == nil || *n != 6 {
		t.Fatalf("mau 6, dapat %v", n)
	}
	n = DaysUntil("2026-09-10", "2026-09-14")
	if n == nil || *n != -4 {
		t.Fatalf("mau -4, dapat %v", n)
	}
	if DaysUntil("rusak", "2026-09-14") != nil {
		t.Fatal("rusak mau nil")
	}
}
