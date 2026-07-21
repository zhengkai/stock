package util

import (
	"fmt"
	"project/config"
	"time"

	"github.com/zhengkai/life-go"
)

type FnSelector func(path string) bool

func IsWorkTime(sleep bool) bool {

	if !config.Prod {
		return true
	}

	hour := time.Now().Hour()
	if hour < 9 || hour > 15 {
		if sleep {
			fmt.Println(`非工作时间`, hour)
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

	weekday := time.Now().Weekday()
	if weekday == time.Saturday || weekday == time.Sunday {
		if sleep {
			fmt.Println(`非工作日`, weekday)
			life.Sleep(32400) // 9 hour
		}
		return false
	}
	return true
}
