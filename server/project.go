// Package project 项目起始
package project

import (
	"project/app"
	"project/config"
	"project/util"
	"project/web"
	"project/zj"
)

func run() {

	zj.Init()

	util.Init()

	if config.Prod {
		go app.Run()
	}
	// go app.Test()

	web.Server()
}

func afterRun() {
}
