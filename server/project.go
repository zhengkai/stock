package project

import (
	"project/app"
	"project/util"
	"project/web"
)

func run() {

	util.Init()

	go app.Run()

	go web.Server()
}

func afterRun() {
}
