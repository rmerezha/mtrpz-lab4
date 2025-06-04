package listener

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/rmerezha/mtrpz-lab4/config"
	"github.com/rmerezha/mtrpz-lab4/runner"
)

type StateWatcherListener struct {
	MasterURL    string
	Host         string
	Runner       runner.Runner
	mu           sync.Mutex
	store        *ContainerStateStore
	pollInterval time.Duration
	Token        string
}

func NewStateWatcherListener(masterURL, host string, r runner.Runner, interval time.Duration, token string) *StateWatcherListener {
	return &StateWatcherListener{
		MasterURL:    masterURL,
		Host:         host,
		Runner:       r,
		store:        NewContainerStateStore(),
		pollInterval: interval,
		Token:        token,
	}
}

func (sw *StateWatcherListener) Listen(stopCh <-chan struct{}) {
	ticker := time.NewTicker(sw.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			sw.checkAndReport()
		}
	}
}

func (sw *StateWatcherListener) checkAndReport() {
	sw.mu.Lock()
	containerNames := make([]string, 0, len(sw.store.states))
	for name := range sw.store.states {
		containerNames = append(containerNames, name)
	}
	sw.mu.Unlock()

	for _, name := range containerNames {
		stateStr, err := sw.Runner.State(name)
		if err != nil {
			log.Printf("StateWatcherListener: failed to get state for %s: %v", name, err)
			continue
		}

		state := config.ContainerState(stateStr)

		sw.mu.Lock()
		prevState, known := sw.store.Get(name)
		if !known || prevState != state {
			sw.store.Set(name, state)
			sw.mu.Unlock()

			sw.sendStateUpdate(name, state)
		} else {
			sw.mu.Unlock()
		}
	}
}

func (sw *StateWatcherListener) sendStateUpdate(containerName string, state config.ContainerState) {
	body := struct {
		Host          string                `json:"host"`
		ContainerName string                `json:"name"`
		State         config.ContainerState `json:"state"`
	}{
		Host:          sw.Host,
		ContainerName: containerName,
		State:         state,
	}

	data, err := json.Marshal(body)
	if err != nil {
		log.Printf("StateWatcherListener: failed to marshal state update: %v", err)
		return
	}

	url := sw.MasterURL + "/api/v1/state"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("StateWatcherListener: failed to create request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+sw.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("StateWatcherListener: failed to send state update: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		log.Printf("StateWatcherListener: unexpected response code %d when sending state update", resp.StatusCode)
		return
	}

	log.Printf("StateWatcherListener: sent state update for %s: %s", containerName, state)
}
