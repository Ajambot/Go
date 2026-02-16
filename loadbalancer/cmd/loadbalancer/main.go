package main

import (
	"loadbalancer/pkg/loadbalancer"
)

func main() {
	lb := loadbalancer.MakeLB("rr")
	lb.Start()
}
