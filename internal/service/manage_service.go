package service

// EstimateMonths: sisa target / kontribusi per bulan, bulat ke atas.
// monthly<=0 atau sudah lunas → -1 (tak terestimasi).
func EstimateMonths(target, current, monthly float64) int {
	rest := target - current
	if rest <= 0 || monthly <= 0 {
		return -1
	}
	m := int(rest / monthly)
	if rest != float64(m)*monthly {
		m++
	}
	return m
}

// FollowupDue: butuh follow-up bila tak pernah kontak, atau
// (today - last) >= followupDays. Tanggal kosong → true (kecuali followupDays<=0).
func FollowupDue(last, today string, followupDays int) bool {
	if followupDays <= 0 {
		return false
	}
	if last == "" || today == "" {
		return true
	}
	d := DaysUntil(last, today)
	if d == nil {
		return true
	}
	return -*d >= followupDays
}

// DocExpiring: kedaluwarsa dalam reminderDays (termasuk sudah lewat).
func DocExpiring(expiry, today string, reminderDays int) bool {
	d := DaysUntil(expiry, today)
	if d == nil {
		return false
	}
	return *d <= reminderDays
}
