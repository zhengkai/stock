package util

import "fmt"

func StockRegion(code string) (string, error) {

	if len(code) != 6 {
		return `unknown`, fmt.Errorf(`invalid stock code length: %d`, len(code))
	}

	f := code[0]

	switch f {
	case '0', '3':
		return `sz`, nil
	case '6':
		return `sh`, nil
	case '8':
		return `bj`, nil
	default:
		return `unknown`, fmt.Errorf(`invalid stock code prefix: %c`, f)
	}
}
