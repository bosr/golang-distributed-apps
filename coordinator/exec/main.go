package main

import (
	"fmt"

	"github.com/bosr/golang-distributed-apps/coordinator"
)

var (
	dc *coordinator.DatabaseConsumer
)

func main() {
	ea := coordinator.NewEventAggregator()
	dc = coordinator.NewDatabaseConsumer(ea)
	ql := coordinator.NewQueueListener(ea)
	go ql.ListenForNewSource()

	var a string
	fmt.Scanln(&a)
}
