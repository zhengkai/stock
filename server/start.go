package project

import (
	"project/config"

	"github.com/zhengkai/life-go"
)

// Start ...
func Start() {

	go run()

	life.Wait()

	afterRun()
}

// Prod ...
func Prod() {

	config.Prod = true

	Start()
}
