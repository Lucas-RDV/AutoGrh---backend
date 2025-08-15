package timeToDateString

import "time"

func TimeToDateString(t time.Time) string {
	return t.Format("2006-01-02")
}
