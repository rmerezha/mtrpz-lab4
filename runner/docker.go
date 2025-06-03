package runner

import (
	"github.com/docker/docker/client"
	"github.com/rmerezha/mtrpz-lab4/config"
)

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
	// TODO
	return nil
}

func (d *DockerRunner) Kill(name string) error {
	// TODO
	return nil
}

func (d *DockerRunner) Restart(name string) error {
	// TODO
	return nil
}

func (d *DockerRunner) Remove(name string) error {
	// TODO
	return nil
}

func (d *DockerRunner) PullImage(name string) error {
	// TODO
	return nil
}
