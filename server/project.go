package project

import (
	"project/app"
	"project/util"
	"project/web"
)

func run() {

	util.Init()

	go app.Run()

	web.Server()
}

func afterRun() {
}
