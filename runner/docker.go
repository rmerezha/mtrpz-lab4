package runner

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/docker/go-units"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/rmerezha/mtrpz-lab4/config"
	"io"
	"strings"
)

const SIGKILL = "SIGKILL"

type DockerClient interface {
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error)
	ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error
	ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerKill(ctx context.Context, containerID, signal string) error
	ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error
	ContainerRemove(ctx context.Context, containerID string, options container.RemoveOptions) error

	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
	ImagePull(ctx context.Context, refStr string, options image.PullOptions) (io.ReadCloser, error)
}

type DockerRunner struct {
	cli DockerClient
}

func NewDockerRunner() (*DockerRunner, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	return &DockerRunner{cli: cli}, nil
}

func (d *DockerRunner) Run(c config.Container) error {
	ctx := context.Background()

	portBindings := nat.PortMap{}
	exposedPorts := nat.PortSet{}
	for _, port := range c.Ports {
		parts := strings.Split(port, ":")
		if len(parts) != 2 {
			return fmt.Errorf("invalid port format: %s", port)
		}
		hostPort := parts[0]
		containerPort := parts[1]

		p, err := nat.NewPort("tcp", containerPort)
		if err != nil {
			return err
		}
		portBindings[p] = []nat.PortBinding{{HostPort: hostPort}}
		exposedPorts[p] = struct{}{}
	}

	cfg := &container.Config{
		Image:        c.Image,
		Env:          toEnvList(c.Environment),
		ExposedPorts: exposedPorts,
	}

	if c.Entrypoint != "" {
		cfg.Entrypoint = strings.Fields(c.Entrypoint)
	}
	if c.Cmd != "" {
		cfg.Cmd = strings.Fields(c.Cmd)
	}

	hostCfg := &container.HostConfig{
		PortBindings: portBindings,
	}

	if err := applyOptions(c.Options, hostCfg, cfg); err != nil {
		return err
	}

	resp, err := d.cli.ContainerCreate(ctx, cfg, hostCfg, nil, nil, c.Name)
	if err != nil {
		return err
	}

	return d.cli.ContainerStart(ctx, resp.ID, container.StartOptions{})
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

func toEnvList(env map[string]string) []string {
	var res []string
	for k, v := range env {
		res = append(res, fmt.Sprintf("%s=%s", k, v))
	}
	return res
}

func applyOptions(opts []string, hostCfg *container.HostConfig, config *container.Config) error {
	for _, opt := range opts {
		switch {
		case strings.HasPrefix(opt, "--net="), strings.HasPrefix(opt, "--network="):
			val := strings.TrimPrefix(strings.TrimPrefix(opt, "--net="), "--network=")
			hostCfg.NetworkMode = container.NetworkMode(val)

		case strings.HasPrefix(opt, "--restart="):
			hostCfg.RestartPolicy = container.RestartPolicy{Name: container.RestartPolicyMode(strings.TrimPrefix(opt, "--restart="))}

		case strings.HasPrefix(opt, "-v "), strings.HasPrefix(opt, "--volume="):
			val := strings.TrimPrefix(strings.TrimPrefix(opt, "--volume="), "-v ")
			hostCfg.Binds = append(hostCfg.Binds, val)

		case opt == "--privileged":
			hostCfg.Privileged = true

		case strings.HasPrefix(opt, "--memory="):
			val := strings.TrimPrefix(opt, "--memory=")
			var mem int64
			_, err := fmt.Sscanf(val, "%d", &mem)
			if err != nil {
				return fmt.Errorf("invalid memory value: %s", val)
			}
			hostCfg.Memory = mem

		case strings.HasPrefix(opt, "--cpus="):
			val := strings.TrimPrefix(opt, "--cpus=")
			var cpus float64
			_, err := fmt.Sscanf(val, "%f", &cpus)
			if err != nil {
				return fmt.Errorf("invalid cpus value: %s", val)
			}
			hostCfg.NanoCPUs = int64(cpus * 1e9)

		case strings.HasPrefix(opt, "--add-host="):
			val := strings.TrimPrefix(opt, "--add-host=")
			hostCfg.ExtraHosts = append(hostCfg.ExtraHosts, val)

		case strings.HasPrefix(opt, "--device="):
			val := strings.TrimPrefix(opt, "--device=")
			hostCfg.Devices = append(hostCfg.Devices, container.DeviceMapping{PathOnHost: val, PathInContainer: val})

		case strings.HasPrefix(opt, "--tmpfs="):
			val := strings.TrimPrefix(opt, "--tmpfs=")
			if hostCfg.Tmpfs == nil {
				hostCfg.Tmpfs = make(map[string]string)
			}
			parts := strings.SplitN(val, ":", 2)
			if len(parts) == 2 {
				hostCfg.Tmpfs[parts[0]] = parts[1]
			} else {
				hostCfg.Tmpfs[val] = ""
			}

		case strings.HasPrefix(opt, "--hostname="):
			config.Hostname = strings.TrimPrefix(opt, "--hostname=")

		case strings.HasPrefix(opt, "--cap-add="):
			hostCfg.CapAdd = append(hostCfg.CapAdd, strings.TrimPrefix(opt, "--cap-add="))

		case strings.HasPrefix(opt, "--cap-drop="):
			hostCfg.CapDrop = append(hostCfg.CapDrop, strings.TrimPrefix(opt, "--cap-drop="))

		case strings.HasPrefix(opt, "--security-opt="):
			hostCfg.SecurityOpt = append(hostCfg.SecurityOpt, strings.TrimPrefix(opt, "--security-opt="))

		case strings.HasPrefix(opt, "--ipc="):
			hostCfg.IpcMode = container.IpcMode(strings.TrimPrefix(opt, "--ipc="))

		case strings.HasPrefix(opt, "--shm-size="):
			val := strings.TrimPrefix(opt, "--shm-size=")
			var size int64
			_, err := fmt.Sscanf(val, "%d", &size)
			if err != nil {
				return fmt.Errorf("invalid shm-size: %s", val)
			}
			hostCfg.ShmSize = size

		case strings.HasPrefix(opt, "--ulimit="):
			val := strings.TrimPrefix(opt, "--ulimit=")
			parts := strings.Split(val, "=")
			if len(parts) != 2 {
				return fmt.Errorf("invalid ulimit format: %s", val)
			}
			name := parts[0]
			limits := strings.Split(parts[1], ":")
			if len(limits) != 2 {
				return fmt.Errorf("invalid ulimit range: %s", parts[1])
			}
			var soft, hard int64
			_, err := fmt.Sscanf(limits[0], "%d", &soft)
			if err != nil {
				return fmt.Errorf("invalid soft limit: %s", limits[0])
			}
			_, err = fmt.Sscanf(limits[1], "%d", &hard)
			if err != nil {
				return fmt.Errorf("invalid hard limit: %s", limits[1])
			}
			hostCfg.Ulimits = append(hostCfg.Ulimits, &units.Ulimit{Name: name, Soft: soft, Hard: hard})

		case strings.HasPrefix(opt, "--dns="):
			hostCfg.DNS = append(hostCfg.DNS, strings.TrimPrefix(opt, "--dns="))

		case strings.HasPrefix(opt, "--dns-search="):
			hostCfg.DNSSearch = append(hostCfg.DNSSearch, strings.TrimPrefix(opt, "--dns-search="))

		case strings.HasPrefix(opt, "--label="):
			val := strings.TrimPrefix(opt, "--label=")
			if config.Labels == nil {
				config.Labels = make(map[string]string)
			}
			kv := strings.SplitN(val, "=", 2)
			if len(kv) == 2 {
				config.Labels[kv[0]] = kv[1]
			} else {
				config.Labels[kv[0]] = ""
			}

		default:
			return fmt.Errorf("unsupported or unknown option: %s", opt)
		}
	}
	return nil
}
