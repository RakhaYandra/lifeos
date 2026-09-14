package service

import "testing"

func TestGoalProgress(t *testing.T) {
	if GoalProgress(100, 25) != 25 {
		t.Fatal("mau 25")
	}
	if GoalProgress(0, 0) != 0 {
		t.Fatal("target 0 + current 0 mau 0")
	}
	if GoalProgress(0, 5) != 100 {
		t.Fatal("target 0 + current>0 mau 100")
	}
	if GoalProgress(10, 99) != 100 {
		t.Fatal("mau cap 100")
	}
	if GoalProgress(10, -5) != 0 {
		t.Fatal("negatif mau 0")
	}
}

func TestDailyStreak(t *testing.T) {
	done := map[string]bool{"2026-09-12": true, "2026-09-13": true, "2026-09-14": true}
	if DailyStreak(done, "2026-09-14") != 3 {
		t.Fatal("mau 3")
	}
	// today belum: hitung dari kemarin
	if DailyStreak(done, "2026-09-15") != 3 {
		t.Fatal("mau 3 dari kemarin")
	}
	// bolong sehari
	done2 := map[string]bool{"2026-09-12": true, "2026-09-14": true}
	if DailyStreak(done2, "2026-09-14") != 1 {
		t.Fatal("bolong mau reset ke 1")
	}
	if DailyStreak(map[string]bool{}, "2026-09-14") != 0 {
		t.Fatal("kosong mau 0")
	}
}

func TestWeeklyStreak(t *testing.T) {
	// target 3/pekan: pekan 7-13 Sep full, pekan berjalan 14 Sep baru 1 → streak 1 (pekan lalu)
	done := []string{"2026-09-07", "2026-09-08", "2026-09-09", "2026-09-10", "2026-09-14"}
	if WeeklyStreak(done, 3, "2026-09-14") != 1 {
		t.Fatalf("mau 1, dapat %d", WeeklyStreak(done, 3, "2026-09-14"))
	}
	done2 := []string{"2026-09-07", "2026-09-08", "2026-09-09", "2026-09-14", "2026-09-15", "2026-09-16"}
	if WeeklyStreak(done2, 3, "2026-09-16") != 2 {
		t.Fatalf("mau 2, dapat %d", WeeklyStreak(done2, 3, "2026-09-16"))
	}
	if WeeklyStreak(nil, 2, "2026-09-14") != 0 {
		t.Fatal("kosong mau 0")
	}
}
