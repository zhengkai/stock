package util

import "fmt"

func Price(v int64) string {
	sign := ``
	if v < 0 {
		sign = `-`
		v = -v
	}
	return fmt.Sprintf(`%s%d.%02d`, sign, v/100, v%100)
}
