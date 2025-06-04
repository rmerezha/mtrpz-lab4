package listener

import (
	"sync"

	"github.com/rmerezha/mtrpz-lab4/config"
)

type ContainerStateStore struct {
	mu     sync.RWMutex
	states map[string]config.ContainerState
}

func NewContainerStateStore() *ContainerStateStore {
	return &ContainerStateStore{
		states: make(map[string]config.ContainerState),
	}
}

func (s *ContainerStateStore) Get(name string) (config.ContainerState, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	state, ok := s.states[name]
	return state, ok
}

func (s *ContainerStateStore) Set(name string, state config.ContainerState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[name] = state
}
