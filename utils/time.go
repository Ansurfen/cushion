package utils

import (
	"strconv"
	"time"
)

func FmtTimestamp(ts int64) string {
	return time.Unix(ts, 0).Format("2006-01-02")
}

func NowTimestamp() int64 {
	return time.Now().Unix()
}

func NowTimestampByString() string {
	return strconv.Itoa(int(time.Now().Unix()))
}
