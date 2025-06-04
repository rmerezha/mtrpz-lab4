package planner

import (
	"github.com/rmerezha/mtrpz-lab4/config"
	"sync"
)

type ContainerState = string

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
	Config       config.Container
	State        ContainerState
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

func (p *Planner) UpdateState(host, containerName string, newState ContainerState) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	containers, ok := p.storage[host]
	if !ok {
		return false
	}

	for _, cs := range containers {
		if cs.Config.Name == containerName {
			cs.State = newState
			return true
		}
	}
	return false
}

func (p *Planner) ListContainersByHost(host string) []*ContainerStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return append([]*ContainerStatus(nil), p.storage[host]...)
}

func (p *Planner) AddManifest(m *config.Manifest) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, c := range m.Containers {
		cs := &ContainerStatus{
			ManifestName: m.Name,
			Config:       c,
			State:        StateCreated,
		}
		p.storage[c.Host] = append(p.storage[c.Host], cs)
	}
}
