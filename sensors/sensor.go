package main

import (
	"flag"
)

var (
	name     = flag.String("name", "sensor", "name of the sensor")
	freq     = flag.Uint("freq", 5, "update frequency in cycles/sec")
	max      = flag.Float64("max", 5., "maximum value for generated readings")
	min      = flag.Float64("min", 1., "minimum value for generated readings")
	stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")
)

func main() {
	flag.Parse()
}
