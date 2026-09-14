package service

import "sort"

var validTripStatus = map[string]bool{
	"planning": true, "booked": true, "ongoing": true, "done": true, "cancelled": true,
}

func ValidTripStatus(s string) bool { return validTripStatus[s] }

// RankOptions: rank 1-based dari total terbobot, urut turun. Total sama → rank sama.
func RankOptions(totals map[int64]float64) map[int64]int {
	ids := make([]int64, 0, len(totals))
	for id := range totals {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return totals[ids[i]] > totals[ids[j]] })
	ranks := map[int64]int{}
	rank := 0
	var prev float64
	first := true
	for _, id := range ids {
		if first || totals[id] != prev {
			rank++
			prev = totals[id]
			first = false
		}
		ranks[id] = rank
	}
	return ranks
}
