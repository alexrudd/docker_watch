package collector

import (
	"flag"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"golang.org/x/net/context"

	"github.com/prometheus/common/log"
)

var (
	dockerAddr = flag.String("collector.docker.addr", "unix:///var/run/docker.sock", "The location of the docker daemon socket or endpoint")
)

var dc *client.Client

func getDockerClient() (dockerClient *client.Client, err error) {
	if dc == nil {
		log.Debugf("Creating new Docker api client")
		dockerClient, err = client.NewClient(*dockerAddr, "v1.22", nil, nil)
		dc = dockerClient
	}
	return dc, err
}

func getContainerList() ([]types.Container, error) {
	// Update - checks for new/departed containers and scrapes them
	log.Debugf("Fetching list of locally running containers")
	cli, err := getDockerClient()
	if err != nil {
		log.Errorf("Failed to create Docker api client: %s", err.Error())
		return nil, err
	}

	options := types.ContainerListOptions{All: true, Quiet: true}
	containers, err := cli.ContainerList(context.Background(), options)
	if err != nil {
		log.Errorf("Failed to fetch container list: %s", err.Error())
	}

	return containers, err
}
