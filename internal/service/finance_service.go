package service

import "time"

var validTrxType = map[string]bool{"income": true, "expense": true, "transfer": true}
var validSubFreq = map[string]bool{"weekly": true, "monthly": true, "yearly": true}

func ValidTrxType(s string) bool { return validTrxType[s] }
func ValidSubFreq(s string) bool { return validSubFreq[s] }

// Utilization: actual/budget*100. budget<=0 → 0.
func Utilization(actual, budget float64) float64 {
	if budget <= 0 {
		return 0
	}
	return actual * 100 / budget
}

// BudgetStatus: safe (<warn), warning (<100), over (>=100).
func BudgetStatus(utilPct, warnPct float64) string {
	if utilPct >= 100 {
		return "over"
	}
	if utilPct >= warnPct {
		return "warning"
	}
	return "safe"
}

// AnnualCost: estimasi biaya tahunan langganan.
func AnnualCost(cost float64, freq string) float64 {
	switch freq {
	case "weekly":
		return cost * 52
	case "yearly":
		return cost
	default:
		return cost * 12
	}
}

// DaysUntil: target - today dalam hari. Nil bila tanggal rusak.
func DaysUntil(target, today string) *int {
	if target == "" || today == "" {
		return nil
	}
	d, err1 := time.Parse("2006-01-02", target)
	t, err2 := time.Parse("2006-01-02", today)
	if err1 != nil || err2 != nil {
		return nil
	}
	n := int(d.Sub(t).Hours() / 24)
	return &n
}
