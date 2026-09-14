package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/RakhaYandra/lifeos/internal/repository"
	"github.com/RakhaYandra/lifeos/internal/service"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	Tasks   *repository.TaskRepository
	Habits  *repository.HabitRepository
	Trx     *repository.TransactionRepository
	Goals   *repository.GoalRepository
	Subs    *repository.SubscriptionRepository
	Reminds *repository.ReminderRepository
}

func (h *DashboardHandler) Get(c *gin.Context) {
	uid, _ := c.Get("userID")
	id := uid.(int64)
	today := service.TodayWIB()
	y, _ := strconv.Atoi(today[:4])
	m, _ := strconv.Atoi(today[5:7])

	allTasks, _ := h.Tasks.List(id, repository.TaskFilter{})
	dueToday, overdue := 0, 0
	for _, t := range allTasks {
		due := ""
		if t.DueDate.Valid {
			due = t.DueDate.String
		}
		if due == today && t.Status != "completed" && t.Status != "cancelled" {
			dueToday++
		}
		if service.IsOverdue(t.Status, due, today) {
			overdue++
		}
	}

	habits, _ := h.Habits.List(id)
	habitsDone := 0
	for _, hb := range habits {
		dates, _ := h.Habits.DoneDates(hb.ID)
		for _, d := range dates {
			if d == today {
				habitsDone++
				break
			}
		}
	}

	income, expense, _ := h.Trx.MonthSummary(id, y, m)
	goals, _ := h.Goals.List(id, "")
	activeGoals, atRisk := 0, 0
	for _, g := range goals {
		if g.Status == "completed" || g.Status == "cancelled" {
			continue
		}
		activeGoals++
		if g.Status == "at_risk" {
			atRisk++
		}
	}

	subs, _ := h.Subs.List(id)
	subSoon := 0
	for _, s := range subs {
		if s.Active != 1 {
			continue
		}
		if d := service.DaysUntil(s.NextBilling, today); d != nil && *d >= 0 && *d <= 14 {
			subSoon++
		}
	}

	rems, _ := h.Reminds.List(id)
	remSoon := []gin.H{}
	for _, m := range rems {
		next := service.NextOccurrence(m.Date, m.Recurrence, today)
		if d := service.DaysUntil(next, today); d != nil && *d >= 0 && *d <= 7 {
			remSoon = append(remSoon, gin.H{"id": m.ID, "title": m.Title, "next": next, "days_until": *d})
		}
	}
	if remSoon == nil {
		remSoon = []gin.H{}
	}

	c.JSON(http.StatusOK, gin.H{
		"today":        today,
		"tasks":        gin.H{"due_today": dueToday, "overdue": overdue, "total": len(allTasks)},
		"habits":       gin.H{"total": len(habits), "done_today": habitsDone},
		"finance":      gin.H{"year": y, "month": m, "income": income, "expense": expense, "net": income - expense},
		"goals":        gin.H{"active": activeGoals, "at_risk": atRisk},
		"subs_due_14d": subSoon,
		"reminders_7d": remSoon,
		"_ts":          time.Now().Unix(),
	})
}
