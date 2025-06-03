package runner

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/rmerezha/mtrpz-lab4/config"
	"io"
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
	ctx := context.Background()
	images, err := d.cli.ImageList(ctx, image.ListOptions{})
	if err != nil {
		return err
	}
	for _, img := range images {
		for _, tag := range img.RepoTags {
			if tag == name {
				return nil
			}
		}
	}

	out, err := d.cli.ImagePull(ctx, name, image.PullOptions{})
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(io.Discard, out)
	if err != nil {
		return err
	}

	return nil
}
