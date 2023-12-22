package time

import (
	"fmt"
	"time"
)

type Time struct {
	Time string
	Date string
}

var dayNames = map[time.Weekday]string{
	time.Monday: "Måndag",
	time.Tuesday: "Tisdag",
	time.Wednesday: "Onsdag",
	time.Thursday: "Torsdag",
	time.Friday: "Fredag",
	time.Saturday: "Lördag",
	time.Sunday: "Söndag",
}

func GetTime() Time {
	now := time.Now().Local()
	timeStr := now.Format("15:04")
	dateStr := fmt.Sprintf("%s %s", dayNames[now.Weekday()], now.Format("2/1"))
	return Time{
		Time: timeStr,
		Date: dateStr,
	}
}
