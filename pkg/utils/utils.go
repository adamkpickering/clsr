package utils

import (
	"time"
)

// Tells the caller whether the two passed times have the same date.
func DatesEqual(time1 time.Time, time2 time.Time) bool {
	time1Year, time1Month, time1Day := time1.Date()
	time2Year, time2Month, time2Day := time2.Date()
	return time1Year == time2Year && time1Month == time2Month && time1Day == time2Day
}
