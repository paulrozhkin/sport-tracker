package utils

import "time"

func TruncateToDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func GetTodayUtc() time.Time {
	timeToday := time.Now().UTC()
	return TruncateToDay(timeToday)
}
