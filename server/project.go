// Package project 项目起始
package project

import (
	"project/app"
	"project/config"
	"project/gold"
	"project/util"
	"project/web"
	"project/zj"

	"github.com/zhengkai/life-go"
)

func run() {

	zj.Init()

	util.Init()

	go func() {
		app.Init()
		if config.Prod {
			life.Sleep(50)
			app.Run()
		}
	}()

	// go app.Test()

	go gold.Loop()

	if false && !config.Prod {
		go func() {
			for {
				zj.J(util.JSONPretty(app.GetPB()))
				zj.J()
				life.Sleep(3)
			}
		}()
	}

	zj.J(util.TimeNumber())

	web.Server()
}

func afterRun() {
}
