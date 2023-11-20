package lang

import "time"

func LocalTime(layout string, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, time.Local)
}
