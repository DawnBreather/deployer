package deployment

import "time"

// TODO: Move to common utilities
func GetTimeNow() time.Time {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc)
}

func GetTimeNowString() string {
	loc, _ := time.LoadLocation("UTC")
	return time.Now().In(loc).String()
}
