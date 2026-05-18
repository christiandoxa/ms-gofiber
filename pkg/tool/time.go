package tool

import "time"

func NowUTC() time.Time {
	return time.Now().UTC()
}

func GetExpiration() time.Duration {
	return EndOfDayExpiration(time.Now())
}

func EndOfDayExpiration(now time.Time) time.Duration {
	targetTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, int(time.Second-time.Nanosecond), now.Location())
	if !now.Before(targetTime) {
		targetTime = targetTime.AddDate(0, 0, 1)
	}
	return targetTime.Sub(now)
}
