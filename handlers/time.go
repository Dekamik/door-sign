package handlers

import (
	"fmt"
	"time"
)

type Time struct {
	Time string
	Date string
}

func GetTime() Time {
	now := time.Now().Local()
	timeStr := now.Format("15:04")
	dateStr := fmt.Sprintf("%s %s", getDayName(now), now.Format("02/01"))
	return Time{
		Time: timeStr,
		Date: dateStr,
	}
}

func getDayName(date time.Time) string {
	wd := date.Weekday()
	switch wd {
	case time.Monday:
		return "Måndag"

	case time.Tuesday:
		return "Tisdag"

	case time.Wednesday:
		return "Onsdag"

	case time.Thursday:
		return "Torsdag"

	case time.Friday:
		return "Fredag"

	case time.Saturday:
		return "Lördag"

	case time.Sunday:
		return "Söndag"

	default:
		return "Fel"
	}
}
