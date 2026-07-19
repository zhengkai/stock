package main

import (
	"server/app"
	"server/web"
)

func main() {
	go app.Run()

	web.Server()
}
