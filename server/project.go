package project

import (
	"project/app"
	"project/util"
	"project/web"
	"project/zj"
)

func run() {

	zj.Init()

	util.Init()

	go app.Run()

	web.Server()
}

func afterRun() {
}
