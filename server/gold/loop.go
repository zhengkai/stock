package gold

import (
	"time"

	"github.com/zhengkai/life-go"
)

func Loop() {
	life.Sleep(5)

	for {
		o := &GoldPrice{
			Time: time.Now(),
		}
		o.Price, o.Error = fetch()
		goldPrice = o

		life.Sleep(55)
	}
}
