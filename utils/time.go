package utils

import "time"

func DaysIn(m time.Month, year int) int {
	if m == time.February && IsLeap(year) {
		return 29
	}
	switch m {
	case time.January, time.March, time.May, time.July, time.August, time.October, time.December:
		return 31
	case time.February:
		return 28
	default:
		return 30
	}
}

func IsLeap(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}
