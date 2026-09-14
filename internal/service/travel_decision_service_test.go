package service

import "testing"

func TestRankOptions(t *testing.T) {
	r := RankOptions(map[int64]float64{1: 10, 2: 30, 3: 20})
	if r[2] != 1 || r[3] != 2 || r[1] != 3 {
		t.Fatalf("rank salah: %v", r)
	}
	r = RankOptions(map[int64]float64{1: 10, 2: 10})
	if r[1] != 1 || r[2] != 1 {
		t.Fatalf("seri mau rank sama: %v", r)
	}
	if len(RankOptions(nil)) != 0 {
		t.Fatal("kosong mau kosong")
	}
}
