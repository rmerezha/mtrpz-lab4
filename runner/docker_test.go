package runner

import (
	"bytes"
	"context"
	"errors"
	"github.com/docker/docker/api/types/network"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"io"
	"strings"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/rmerezha/mtrpz-lab4/config"
)

type mockDockerClient struct {
	containerCreated bool
	startCalled      bool
	pulledImages     []string
	existingImages   []string
}

func (m *mockDockerClient) ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error) {
	if containerName == "invalid" {
		return container.CreateResponse{}, errors.New("create failed")
	}
	m.containerCreated = true

	if containerName == "start-err" {
		return container.CreateResponse{}, nil
	}
	return container.CreateResponse{ID: "mocked-container-id"}, nil
}

func (m *mockDockerClient) ContainerStart(ctx context.Context, containerID string, opts container.StartOptions) error {
	if containerID == "mocked-container-id" {
		m.startCalled = true
		return nil
	}
	return errors.New("start failed")
}

func (m *mockDockerClient) ContainerStop(ctx context.Context, id string, opts container.StopOptions) error {
	if id == "test" {
		return nil
	}
	return errors.New("container not found")
}

func (m *mockDockerClient) ContainerKill(ctx context.Context, id, signal string) error {
	if id == "test" && signal == SIGKILL {
		return nil
	}
	return errors.New("kill failed")
}

func (m *mockDockerClient) ContainerRestart(ctx context.Context, id string, opts container.StopOptions) error {
	if id == "test" {
		return nil
	}
	return errors.New("restart failed")
}

func (m *mockDockerClient) ContainerRemove(ctx context.Context, id string, opts container.RemoveOptions) error {
	if id == "test" && opts.Force {
		return nil
	}
	return errors.New("remove failed")
}

func (m *mockDockerClient) ImageList(ctx context.Context, opts image.ListOptions) ([]image.Summary, error) {
	var summaries []image.Summary
	for _, tag := range m.existingImages {
		summaries = append(summaries, image.Summary{RepoTags: []string{tag}})
	}
	return summaries, nil
}

func (m *mockDockerClient) ImagePull(ctx context.Context, ref string, opts image.PullOptions) (io.ReadCloser, error) {
	m.pulledImages = append(m.pulledImages, ref)
	return io.NopCloser(bytes.NewBufferString("pulled")), nil
}

func (m *mockDockerClient) ContainerInspect(ctx context.Context, containerID string) (container.InspectResponse, error) {
	if containerID == "invalid" {
		return container.InspectResponse{}, errors.New("inspect failed")
	}
	return container.InspectResponse{
		ContainerJSONBase: &container.ContainerJSONBase{
			State: &container.State{
				Status: container.StateRunning,
			},
		},
	}, nil
}

func TestDockerRunner_Run(t *testing.T) {
	mock := &mockDockerClient{}
	runner := &DockerRunner{cli: mock}

	containerConfig := config.Container{
		Name:  "test-container",
		Image: "alpine",
		Ports: []string{"8080:80"},
	}

	err := runner.Run(containerConfig)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !mock.containerCreated || !mock.startCalled {
		t.Error("expected ContainerCreate and ContainerStart to be called")
	}
}

func TestDockerRunner_Run_InvalidPortFormat(t *testing.T) {
	mock := &mockDockerClient{}
	runner := &DockerRunner{cli: mock}

	c := config.Container{
		Name:  "test-container",
		Image: "alpine",
		Ports: []string{"badformat"},
	}

	err := runner.Run(c)
	if err == nil || !strings.Contains(err.Error(), "invalid port format") {
		t.Errorf("expected invalid port format error, got %v", err)
	}
}

func TestDockerRunner_Run_CreateError(t *testing.T) {
	mock := &mockDockerClient{}
	runner := &DockerRunner{cli: mock}

	c := config.Container{
		Name:  "invalid",
		Image: "alpine",
		Ports: []string{"8080:80"},
	}

	err := runner.Run(c)
	if err == nil || !strings.Contains(err.Error(), "create failed") {
		t.Errorf("expected create error, got %v", err)
	}
}

func TestDockerRunner_Run_StartError(t *testing.T) {
	mock := &mockDockerClient{}
	runner := &DockerRunner{cli: mock}

	c := config.Container{
		Name:  "start-err",
		Image: "alpine",
		Ports: []string{"8080:80"},
	}

	err := runner.Run(c)
	if err == nil || !strings.Contains(err.Error(), "start failed") {
		t.Errorf("expected start error, got %v", err)
	}
}

func TestDockerRunner_PullImage_AlreadyExists(t *testing.T) {
	mock := &mockDockerClient{existingImages: []string{"alpine"}}
	runner := &DockerRunner{cli: mock}

	err := runner.PullImage("alpine")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(mock.pulledImages) != 0 {
		t.Error("expected image not to be pulled")
	}
}

func TestDockerRunner_PullImage_NewImage(t *testing.T) {
	mock := &mockDockerClient{}
	runner := &DockerRunner{cli: mock}

	err := runner.PullImage("busybox")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(mock.pulledImages) != 1 || mock.pulledImages[0] != "busybox" {
		t.Error("expected busybox to be pulled")
	}
}

func TestDockerRunner_Stop(t *testing.T) {
	mock := &mockDockerClient{}
	runner := &DockerRunner{cli: mock}

	err := runner.Stop("test")
	if err != nil {
		t.Errorf("expected stop to succeed, got error: %v", err)
	}
}

func TestDockerRunner_Kill(t *testing.T) {
	mock := &mockDockerClient{}
	runner := &DockerRunner{cli: mock}

	err := runner.Kill("test")
	if err != nil {
		t.Errorf("expected kill to succeed, got error: %v", err)
	}
}

func TestDockerRunner_Restart(t *testing.T) {
	mock := &mockDockerClient{}
	runner := &DockerRunner{cli: mock}

	err := runner.Restart("test")
	if err != nil {
		t.Errorf("expected restart to succeed, got error: %v", err)
	}
}

func TestDockerRunner_Remove(t *testing.T) {
	mock := &mockDockerClient{}
	runner := &DockerRunner{cli: mock}

	err := runner.Remove("test")
	if err != nil {
		t.Errorf("expected remove to succeed, got error: %v", err)
	}
}

func TestDockerRunner_State(t *testing.T) {
	mock := &mockDockerClient{}
	runner := &DockerRunner{cli: mock}

	state, err := runner.State("test")
	if err != nil {
		t.Errorf("expected state to succeed, got error: %v", err)
	}
	if state != container.StateRunning {
		t.Errorf("expected state to be running, got %s", state)
	}
}

func TestDockerRunner_State_Fail(t *testing.T) {
	mock := &mockDockerClient{}
	runner := &DockerRunner{cli: mock}

	_, err := runner.State("invalid")
	if err == nil {
		t.Errorf("expected state to fail")
	}
}
