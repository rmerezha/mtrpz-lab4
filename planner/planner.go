package planner

import (
	"github.com/rmerezha/mtrpz-lab4/config"
	"sync"
)

type Planner struct {
	mu      sync.RWMutex
	storage map[string][]*config.ContainerStatus
}

func NewPlanner(manifests ...*config.Manifest) *Planner {
	p := &Planner{
		storage: make(map[string][]*config.ContainerStatus),
	}

	for _, m := range manifests {
		for _, c := range m.Containers {
			cs := &config.ContainerStatus{
				ManifestName: m.Name,
				Config:       c,
				State:        config.StateCreated,
			}
			p.storage[c.Host] = append(p.storage[c.Host], cs)
		}
	}
	return p
}

func (p *Planner) UpdateState(host, containerName string, newState config.ContainerState) bool {
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

func (p *Planner) ListContainersByHost(host string) []*config.ContainerStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return append([]*config.ContainerStatus(nil), p.storage[host]...)
}

func (p *Planner) AddManifest(m *config.Manifest) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, c := range m.Containers {
		cs := &config.ContainerStatus{
			ManifestName: m.Name,
			Config:       c,
			State:        config.StateNew,
		}
		p.storage[c.Host] = append(p.storage[c.Host], cs)
	}
}

func (p *Planner) MarkManifestRemoving(name string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()

	found := false

	for _, containers := range p.storage {
		for _, cs := range containers {
			if cs.ManifestName == name {
				cs.State = config.StateRemoving
				found = true
			}
		}
	}

	return found
}

func (p *Planner) ListContainersByManifest(name string) []*config.ContainerStatus {
	p.mu.RLock()
	defer p.mu.RUnlock()

	var result []*config.ContainerStatus
	for _, containers := range p.storage {
		for _, cs := range containers {
			if name == "" || cs.ManifestName == name {
				result = append(result, cs)
			}
		}
	}
	return result
}
