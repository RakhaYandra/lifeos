package service

import (
	"regexp"
	"time"
)

var monthRe = regexp.MustCompile(`^\d{4}-(0[1-9]|1[0-2])$`)
var yearRe = regexp.MustCompile(`^\d{4}$`)

func ValidMonthPeriod(s string) bool { return monthRe.MatchString(s) }
func ValidYearPeriod(s string) bool  { return yearRe.MatchString(s) }

// MonthBounds: YYYY-MM → tanggal awal/akhir bulan.
func MonthBounds(period string) (from, to string, ok bool) {
	if !ValidMonthPeriod(period) {
		return "", "", false
	}
	t, err := time.Parse("2006-01", period)
	if err != nil {
		return "", "", false
	}
	end := t.AddDate(0, 1, -1)
	return t.Format("2006-01-02"), end.Format("2006-01-02"), true
}

// YearBounds: YYYY → 1 Jan s.d. 31 Des.
func YearBounds(period string) (from, to string, ok bool) {
	if !ValidYearPeriod(period) {
		return "", "", false
	}
	return period + "-01-01", period + "-12-31", true
}

// CountInRange: hitung tanggal YYYY-MM-DD dalam [from,to].
func CountInRange(dates []string, from, to string) int {
	n := 0
	for _, d := range dates {
		if d >= from && d <= to {
			n++
		}
	}
	return n
}
