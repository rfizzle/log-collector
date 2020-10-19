package akamai

import (
	"strconv"
	"time"
)

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func UnixStringToTime(s string) time.Time {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return time.Unix(i, 0)
}
