package app

import (
	"fmt"
	"project/tc"
	"project/util"
	"time"
)

func (a *App) checkAlert() {
	for _, row := range a.Config.GetAlert() {
		s, err := tc.Stock(row.GetCode())
		fmt.Println(util.JSON(s), err)
		time.Sleep(time.Second)
	}
}
