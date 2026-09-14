package service

import "time"

var validGoalLevel = map[string]bool{"annual": true, "quarterly": true, "monthly": true}

var validGoalStatus = map[string]bool{
	"not_started": true, "active": true, "on_track": true,
	"at_risk": true, "completed": true, "cancelled": true,
}

var validHabitFreq = map[string]bool{"daily": true, "weekly": true}

func ValidGoalLevel(s string) bool  { return validGoalLevel[s] }
func ValidGoalStatus(s string) bool { return validGoalStatus[s] }
func ValidHabitFreq(s string) bool  { return validHabitFreq[s] }

// ValidGoalParent: cascade annual→quarterly→monthly.
// quarterly wajib berparent annual; monthly boleh quarterly/annual;
// annual tak berparent. parentLevel "" = tanpa parent.
func ValidGoalParent(level, parentLevel string) bool {
	switch level {
	case "annual":
		return parentLevel == ""
	case "quarterly":
		return parentLevel == "annual"
	case "monthly":
		return parentLevel == "" || parentLevel == "quarterly" || parentLevel == "annual"
	default:
		return false
	}
}

// GoalProgress: current/target*100, cap 0-100. target<=0 → 0 bila current<=0, else 100 bila current>0.
func GoalProgress(target, current float64) int {
	if target <= 0 {
		if current > 0 {
			return 100
		}
		return 0
	}
	p := int(current * 100 / target)
	if p < 0 {
		return 0
	}
	if p > 100 {
		return 100
	}
	return p
}

// DailyStreak: hari berurutan done mundur dari today. Bila today belum done, mulai dari kemarin.
func DailyStreak(done map[string]bool, today string) int {
	t, err := time.Parse("2006-01-02", today)
	if err != nil {
		return 0
	}
	if !done[t.Format("2006-01-02")] {
		t = t.AddDate(0, 0, -1)
	}
	n := 0
	for {
		if !done[t.Format("2006-01-02")] {
			return n
		}
		n++
		t = t.AddDate(0, 0, -1)
	}
}

// WeeklyStreak: pekan berurutan (Senin-Minggu) dengan doneCount>=target, mundur dari pekan today.
// Bila pekan berjalan belum capai target, mulai dari pekan lalu.
func WeeklyStreak(doneDates []string, target int, today string) int {
	if target <= 0 {
		target = 1
	}
	t, err := time.Parse("2006-01-02", today)
	if err != nil {
		return 0
	}
	counts := map[string]int{}
	for _, d := range doneDates {
		p, err := time.Parse("2006-01-02", d)
		if err != nil {
			continue
		}
		wd := int(p.Weekday())
		if wd == 0 {
			wd = 7
		}
		mon := p.AddDate(0, 0, -(wd - 1)).Format("2006-01-02")
		counts[mon]++
	}
	wd := int(t.Weekday())
	if wd == 0 {
		wd = 7
	}
	mon := t.AddDate(0, 0, -(wd - 1))
	if counts[mon.Format("2006-01-02")] < target {
		mon = mon.AddDate(0, 0, -7)
	}
	n := 0
	for {
		if counts[mon.Format("2006-01-02")] < target {
			return n
		}
		n++
		mon = mon.AddDate(0, 0, -7)
	}
}
