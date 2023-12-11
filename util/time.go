package util

import "time"

func UnixYyyyMmDd() string {
	return time.Unix(time.Now().Unix(), 0).Format("2006-01-02")
}

func UnixSeconds() int64 {
	return time.Now().Unix()
}
