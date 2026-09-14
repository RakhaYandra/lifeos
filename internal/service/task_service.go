package service

import "time"

var wib = time.FixedZone("WIB", 7*3600)

var validTaskStatus = map[string]bool{
	"inbox": true, "not_started": true, "in_progress": true,
	"waiting": true, "completed": true, "cancelled": true,
}

var validTaskPriority = map[string]bool{
	"low": true, "medium": true, "high": true, "critical": true,
}

var validProjectStatus = map[string]bool{
	"planning": true, "active": true, "on_hold": true, "completed": true, "cancelled": true,
}

func ValidTaskStatus(s string) bool    { return validTaskStatus[s] }
func ValidTaskPriority(p string) bool  { return validTaskPriority[p] }
func ValidProjectStatus(s string) bool { return validProjectStatus[s] }

func TodayWIB() string { return time.Now().In(wib).Format("2006-01-02") }

// IsOverdue: due < today dan belum selesai/batal. Tanpa due → false.
func IsOverdue(status, due, today string) bool {
	if due == "" || today == "" {
		return false
	}
	if status == "completed" || status == "cancelled" {
		return false
	}
	return due < today
}

// DaysRemaining: due - today dalam hari. Nil bila tanpa due.
func DaysRemaining(due, today string) *int {
	if due == "" || today == "" {
		return nil
	}
	d, err1 := time.Parse("2006-01-02", due)
	t, err2 := time.Parse("2006-01-02", today)
	if err1 != nil || err2 != nil {
		return nil
	}
	n := int(d.Sub(t).Hours() / 24)
	return &n
}

// ProjectProgress: done/total*100, cap 0-100, total 0 → 0.
func ProjectProgress(total, done int) int {
	if total <= 0 || done <= 0 {
		return 0
	}
	p := done * 100 / total
	if p > 100 {
		return 100
	}
	return p
}
