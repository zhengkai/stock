package util

import "time"

func TS() uint32 {
	return uint32(time.Now().Unix())
}
