package docker

import (
	"context"
	"fmt"
	"log"
	"net/netip"
	"sync"

	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/api/types/network"
	"github.com/moby/moby/client"
)

var (
	instance *docker
	once     sync.Once
)

type docker struct {
	client *client.Client
}

func GetInstance() (*docker, error) {
	var err error
	once.Do(func() {
		instance = &docker{}
		err = instance.createClient()
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("Instantiated docker client.")
	return instance, nil
}

func (d *docker) createClient() error {
	cli, err := client.New(
		client.FromEnv,
	)
	if err != nil {
		return err
	}
	d.client = cli
	return nil
}

func (d *docker) CreateContainer(port int) (string, error) {
	if d.client == nil {
		log.Fatal("Error: cannot find Docker client")
	}

	ctx := context.Background()

	netPort, err := network.ParsePort(fmt.Sprintf("%d/tcp", port))
	if err != nil {
		return "", err
	}

	containerPort := network.Port(netPort)
	hostip, err := netip.ParseAddr("0.0.0.0")
	if err != nil {
		return "", err
	}

	resp, err := d.client.ContainerCreate(
		ctx,
		client.ContainerCreateOptions{
			Config: &container.Config{
				ExposedPorts: network.PortSet{
					containerPort: struct{}{},
				},
				Cmd: []string{fmt.Sprintf("%d", port), "1"},
			},
			HostConfig: &container.HostConfig{
				PortBindings: network.PortMap{
					containerPort: []network.PortBinding{
						{
							HostIP:   hostip,
							HostPort: fmt.Sprintf("%d", port),
						},
					},
				},
			},
			Image: "httpserver",
		},
	)
	if err != nil {
		return "", err
	}

	if _, err := d.client.ContainerStart(ctx, resp.ID, client.ContainerStartOptions{}); err != nil {
		return "", err
	}
	return resp.ID, nil
}

func (d *docker) StopContainer(id string) error {
	if d.client == nil {
		log.Fatal("Error: cannot find Docker client")
	}

	ctx := context.Background()

	_, err := d.client.ContainerStop(ctx, id, client.ContainerStopOptions{})

	if err != nil {
		return err
	}
	return nil
}
