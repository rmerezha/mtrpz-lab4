package runner

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/rmerezha/mtrpz-lab4/config"
)

const SIGKILL = "SIGKILL"

type DockerRunner struct {
	cli *client.Client
}

func NewDockerRunner() (*DockerRunner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerRunner{cli: cli}, nil
}

func (d *DockerRunner) Run(c config.Container) error {
	// TODO
	return nil
}

func (d *DockerRunner) Stop(name string) error {
	return d.cli.ContainerStop(context.Background(), name, container.StopOptions{})
}

func (d *DockerRunner) Kill(name string) error {
	return d.cli.ContainerKill(context.Background(), name, SIGKILL)
}

func (d *DockerRunner) Restart(name string) error {
	return d.cli.ContainerRestart(context.Background(), name, container.StopOptions{})
}

func (d *DockerRunner) Remove(name string) error {
	return d.cli.ContainerRemove(context.Background(), name, container.RemoveOptions{Force: true})
}

func (d *DockerRunner) PullImage(name string) error {
	// TODO
	return nil
}
