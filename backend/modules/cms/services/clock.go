package services

import "time"

func timeNowNano() int64 { return time.Now().UnixNano() }
