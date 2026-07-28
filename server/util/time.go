package util

import (
	"strconv"
	"strings"
	"time"
)

func TS() uint32 {
	return uint32(time.Now().Unix())
}

// TimeNumber 让类似 09:30 的时间变成 930，方便直接做数字比较
func TimeNumber() int {

	s := time.Now().Format(`1504`)
	s = strings.TrimLeft(s, `0`)

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}
