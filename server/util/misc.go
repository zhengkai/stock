package util

import (
	"project/config"
	"project/zj"
	"time"

	"github.com/zhengkai/life-go"
)

type FnSelector func(path string) bool

func IsWorkTime(sleep bool) bool {

	if !config.Prod {
		return true
	}

	weekday := time.Now().Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		if sleep {
			zj.J(`非工作日`, weekday)
			life.Sleep(32400) // 9 hour
		}
		return false
	}

	hour := time.Now().Hour()
	if hour < 9 || hour > 15 {
		if sleep {
			zj.J(`非工作时间`, hour)
			if hour == 8 {
				life.Sleep(600)
			} else if hour > 15 {
				life.Sleep(32400)
			} else {
				life.Sleep(3600)
			}
		}
		return false
	}

	return true
}
