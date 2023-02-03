package utils

import (
	"fmt"
	"time"
)

// Tells the caller whether the two passed times have the same date.
func DatesEqual(time1 time.Time, time2 time.Time) bool {
	time1Year, time1Month, time1Day := time1.Date()
	time2Year, time2Month, time2Day := time2.Date()
	return time1Year == time2Year && time1Month == time2Month && time1Day == time2Day
}

// Returns the time in hours (if on the same day) or days (if on different days)
// between two time.Times as a human-readable string. For example, "11h" or "23d".
// time1 should be before time2, otherwise a negative duration will be returned.
func GetReadableTimeDifference(time1, time2 time.Time) string {
	if DatesEqual(time1, time2) {
		difference := time2.Sub(time1)
		hours := difference / time.Hour
		return fmt.Sprintf("%dh", hours)
	} else {
		midnightTime1 := time.Date(time1.Year(), time1.Month(), time1.Day(), 0, 0, 0, 0, time1.Location())
		difference := time2.Sub(midnightTime1)
		day := 24 * time.Hour
		days := difference / day
		return fmt.Sprintf("%dd", days)
	}
}
