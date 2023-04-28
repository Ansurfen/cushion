package utils

import (
	"strconv"
	"time"
)

// FmtTimestamp return a time string whom format is 2006-01-02
func FmtTimestamp(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02")
}

// NowTimestamp returm current timestamp
func NowTimestamp() int64 {
	return time.Now().Unix()
}

// NowTimestamp returm current timestamp string
func NowTimestampByString() string {
	return strconv.Itoa(int(time.Now().Unix()))
}
