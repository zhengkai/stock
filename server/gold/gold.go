package gold

import "time"

type GoldPrice struct {
	Price float64
	Time  time.Time
	Error error
}

var goldPrice = &GoldPrice{}

func Get() *GoldPrice {
	return goldPrice
}
