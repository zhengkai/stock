package app

import (
	"errors"
	"fmt"
	"project/util"
	"time"

	"github.com/zhengkai/life-go"
)

func Run() {
	a := &App{}

	for !life.Stop {

		fmt.Println(`tick`, time.Now().Format(time.DateTime))
		if !util.IsWorkTime(true) {
			continue
		}

		err := a.run()
		if err != nil {
			fmt.Println(`app loop error:`, err)
		}
		life.Sleep(50)
	}
}

func (app *App) run() error {

	err := app.loadConfig()
	if err != nil {
		return err
	}

	app.checkAlert()

	app.logLimit()

	// fmt.Println(tc.StockURL([]string{`600519`}))

	return nil
}

func (app *App) logLimit() error {

	code := `002594`
	q := Stock(code)
	if q == nil {
		return errors.New(`stock query error`)
	}

	f := util.NewFile(`stock/limit.log`)
	return f.AppendF("%s %s %6d %6d\n", time.Now().Format(time.DateTime), code, q.GetLimitUp(), q.GetLimitDown())
}
