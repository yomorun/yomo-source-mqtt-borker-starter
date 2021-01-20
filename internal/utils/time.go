package utils

import "time"

func Now() int64 {
	now := time.Now()
	return now.UnixNano() / 1e6
}
