package service

import "time"

var validIntensity = map[string]bool{"low": true, "medium": true, "high": true}
var validLearnType = map[string]bool{
	"course": true, "book": true, "tutorial": true,
	"certification": true, "practice": true, "other": true,
}
var validLearnStatus = map[string]bool{
	"not_started": true, "active": true, "completed": true, "cancelled": true,
}
var validReadType = map[string]bool{
	"book": true, "article": true, "research": true, "docs": true, "other": true,
}
var validReadStatus = map[string]bool{
	"planned": true, "reading": true, "finished": true, "dropped": true,
}
var validRecurrence = map[string]bool{
	"none": true, "daily": true, "weekly": true, "monthly": true, "yearly": true,
}

func ValidIntensity(s string) bool   { return validIntensity[s] }
func ValidLearnType(s string) bool   { return validLearnType[s] }
func ValidLearnStatus(s string) bool { return validLearnStatus[s] }
func ValidReadType(s string) bool    { return validReadType[s] }
func ValidReadStatus(s string) bool  { return validReadStatus[s] }
func ValidRecurrence(s string) bool  { return validRecurrence[s] }

// WeekRange: Senin-Minggu berisi today (WIB bila today kosong → hari ini).
func WeekRange(today string) (start, end string) {
	t := time.Now().In(wib)
	if today != "" {
		if p, err := time.Parse("2006-01-02", today); err == nil {
			t = p
		}
	}
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	mon := t.AddDate(0, 0, -(wd - 1))
	return mon.Format("2006-01-02"), mon.AddDate(0, 0, 6).Format("2006-01-02")
}

// NextOccurrence: tanggal berikutnya untuk reminder berulang dari today.
func NextOccurrence(date, recurrence, today string) string {
	if recurrence == "none" || date == "" || today == "" {
		return date
	}
	d, err1 := time.Parse("2006-01-02", date)
	t, err2 := time.Parse("2006-01-02", today)
	if err1 != nil || err2 != nil {
		return date
	}
	for !d.After(t) {
		switch recurrence {
		case "daily":
			d = d.AddDate(0, 0, 1)
		case "weekly":
			d = d.AddDate(0, 0, 7)
		case "monthly":
			d = d.AddDate(0, 1, 0)
		case "yearly":
			d = d.AddDate(1, 0, 0)
		default:
			return d.Format("2006-01-02")
		}
	}
	return d.Format("2006-01-02")
}
