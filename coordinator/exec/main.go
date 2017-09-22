package main

import (
	"fmt"

	"github.com/bosr/golang-distributed-apps/coordinator"
)

func main() {
	ql := coordinator.NewQueueListener()
	go ql.ListenForNewSource()

	var a string
	fmt.Scanln(&a)
}
