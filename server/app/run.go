package app

import (
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

func (a *App) run() error {

	err := a.loadConfig()
	if err != nil {
		return err
	}

	a.checkAlert()

	// fmt.Println(tc.StockURL([]string{`600519`}))

	return nil
}
