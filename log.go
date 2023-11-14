package main

import (
	"fmt"
	"time"
)

const TIMESTAMP_FORMAT = "15:04:05.000"

// todo: actual logging
func plog(msg any) {
	timeStr := time.Now().Format(TIMESTAMP_FORMAT)
	data := fmt.Sprintf("%v", msg)
	println(timeStr, data)
}
