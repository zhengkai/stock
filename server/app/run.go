package app

import (
	"fmt"
	"server/tc"
	"server/util"
	"time"
)

func Run() {
	a := &App{}
	a.run()
}

func (a *App) run() error {
	err := a.loadConfig()
	if err != nil {
		return err
	}

	fmt.Println(a.Config)

	for _, row := range a.Config.GetAlert() {
		s, err := tc.Stock(row.GetCode())
		fmt.Println(util.JSON(s), err)
		time.Sleep(time.Second)
	}

	// fmt.Println(tc.StockURL([]string{`600519`}))

	return nil
}
