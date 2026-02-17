package main

import (
	"loadbalancer/pkg/loadbalancer"
	"loadbalancer/pkg/server"
)

func main() {
	lb := loadbalancer.MakeLB("rr")
	lb.Register(server.Server{Load: 0, Healthy: true, Id: 6767})
	lb.Register(server.Server{Load: 0, Healthy: true, Id: 6967})
	lb.Register(server.Server{Load: 0, Healthy: true, Id: 6969})
	lb.Start()
}
