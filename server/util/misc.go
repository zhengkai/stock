package util

import (
	"project/config"
	"project/zj"
	"time"

	"github.com/zhengkai/life-go"
)

const (
	MimeProto     = `application/protobuf`
	MimeProtoJSON = `application/protobuf+json; charset=utf-8`
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

	tn := TimeNumber()
	if tn < 930 || tn > 1500 {
		if sleep {
			zj.J(`非工作时间`, time.Now().Format(time.TimeOnly))
			if tn < 930 {
				life.Sleep(600)
			} else if tn > 1500 {
				life.Sleep(32400)
			} else {
				life.Sleep(3600)
			}
		}
		return false
	}

	return true
}
