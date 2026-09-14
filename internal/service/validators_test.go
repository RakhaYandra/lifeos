package service

import (
	"regexp"
	"testing"
	"time"
)

// Table test validator enum + helper tanggal. Menaikkan mean per-fungsi
// (gate CI: rata-rata % per fungsi internal/service >= 60).
func TestValidators(t *testing.T) {
	cases := []struct {
		name string
		got  bool
		want bool
	}{
		{"task status ok", ValidTaskStatus("inbox"), true},
		{"task status bad", ValidTaskStatus("ngawur"), false},
		{"task priority ok", ValidTaskPriority("high"), true},
		{"task priority bad", ValidTaskPriority("x"), false},
		{"project status ok", ValidProjectStatus("active"), true},
		{"project status bad", ValidProjectStatus("x"), false},
		{"goal level quarterly", ValidGoalLevel("quarterly"), true},
		{"goal level bad", ValidGoalLevel("dekade"), false},
		{"goal status ok", ValidGoalStatus("at_risk"), true},
		{"goal status bad", ValidGoalStatus("x"), false},
		{"habit freq ok", ValidHabitFreq("weekly"), true},
		{"habit freq bad", ValidHabitFreq("hourly"), false},
		{"trx type ok", ValidTrxType("transfer"), true},
		{"trx type bad", ValidTrxType("x"), false},
		{"sub freq ok", ValidSubFreq("yearly"), true},
		{"sub freq bad", ValidSubFreq("daily-x"), false},
		{"intensity ok", ValidIntensity("low"), true},
		{"intensity bad", ValidIntensity("x"), false},
		{"learn type ok", ValidLearnType("book"), true},
		{"learn type bad", ValidLearnType("x"), false},
		{"learn status ok", ValidLearnStatus("completed"), true},
		{"learn status bad", ValidLearnStatus("x"), false},
		{"read type ok", ValidReadType("docs"), true},
		{"read type bad", ValidReadType("x"), false},
		{"read status ok", ValidReadStatus("finished"), true},
		{"read status bad", ValidReadStatus("x"), false},
		{"recurrence ok", ValidRecurrence("monthly"), true},
		{"recurrence bad", ValidRecurrence("x"), false},
		{"trip status ok", ValidTripStatus("done"), true},
		{"trip status bad", ValidTripStatus("x"), false},
		{"month period ok", ValidMonthPeriod("2026-09"), true},
		{"month period bad", ValidMonthPeriod("2026-13"), false},
		{"year period ok", ValidYearPeriod("2026"), true},
		{"year period bad", ValidYearPeriod("26"), false},
	}
	for _, c := range cases {
		if c.got != c.want {
			t.Fatalf("%s: mau %v", c.name, c.want)
		}
	}
	if !ValidGoalParent("annual", "") || !ValidGoalParent("quarterly", "annual") {
		t.Fatal("cascade dasar")
	}
	if ok := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`).MatchString(TodayWIB()); !ok {
		t.Fatalf("TodayWIB format: %s", TodayWIB())
	}
	s, e := WeekRange("")
	if s == "" || e == "" {
		t.Fatal("WeekRange default kosong")
	}
	if _, err := time.Parse("2006-01-02", TodayWIB()); err != nil {
		t.Fatalf("TodayWIB tak parse: %v", err)
	}
}
