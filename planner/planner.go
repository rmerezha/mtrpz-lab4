package planner

import (
	"github.com/rmerezha/mtrpz-lab4/config"
	"sync"
)

type ContainerState = string

const (
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
	Config       config.Container
	State        ContainerState
	Reason       string
}

type Planner struct {
	mu      sync.RWMutex
	storage map[string][]*ContainerStatus
}

func NewPlanner(manifests ...*config.Manifest) *Planner {
	p := &Planner{
		storage: make(map[string][]*ContainerStatus),
	}

	for _, m := range manifests {
		for _, c := range m.Containers {
			cs := &ContainerStatus{
				ManifestName: m.Name,
				Config:       c,
				State:        StateCreated,
			}
			p.storage[c.Host] = append(p.storage[c.Host], cs)
		}
	}
	return p
}
