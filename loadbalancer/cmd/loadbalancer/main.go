package main

import (
	"loadbalancer/pkg/docker"
	"loadbalancer/pkg/loadbalancer"
	"loadbalancer/pkg/server"
)

func main() {
	lb, err := loadbalancer.MakeLB("rr")
	if err != nil {
		panic(err)
	}

	d, err := docker.GetInstance()
	if err != nil {
		panic(err)
	}

	id, err := d.CreateContainer(6767)
	if err != nil {
		panic(err)
	}
	defer d.StopContainer(id)

	id, err = d.CreateContainer(6967)
	if err != nil {
		panic(err)
	}
	defer d.StopContainer(id)

	id, err = d.CreateContainer(6969)
	if err != nil {
		panic(err)
	}
	defer d.StopContainer(id)

	lb.Register(server.Server{Load: 0, Healthy: true, Id: 6767})
	lb.Register(server.Server{Load: 0, Healthy: true, Id: 6967})
	lb.Register(server.Server{Load: 0, Healthy: true, Id: 6969})
	lb.Start()
}
