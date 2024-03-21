package utils

import "time"

// TODO: Check how to check things with timezones

// SleepTime returns the time to sleep until the next update
// It returns a different time depending on the current time
// It is more important to get updates during school hours
// and less important to get updates during the night
// It is also less important to get updates on weekends
// It more regularily checks for updates during the morning when changes usually happen
func SleepTime(t time.Time) time.Duration {
	// Update less often on weekends
	if t.Weekday() == time.Saturday || t.Weekday() == time.Sunday {
		return 40 * time.Minute
	}

	// On Week-days
	if t.Hour() > 12 {
		return 22 * time.Minute
	} else if t.Hour() == 6 {
		return 15 * time.Minute
	} else if t.Hour() > 6 {
		return (7*60 + 30) * time.Second // 7min 30sec
	} else {
		return 33 * time.Minute
	}
}
