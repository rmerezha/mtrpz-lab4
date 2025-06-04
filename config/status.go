package config

type ContainerState string

const (
	StateNew        ContainerState = "new"
	StateCreated    ContainerState = "created"
	StateRunning    ContainerState = "running"
	StatePaused     ContainerState = "paused"
	StateRestarting ContainerState = "restarting"
	StateRemoving   ContainerState = "removing"
	StateExited     ContainerState = "exited"
	StateDead       ContainerState = "dead"
)

type ContainerStatus struct {
	ManifestName string
	Config       Container
	State        ContainerState
}
