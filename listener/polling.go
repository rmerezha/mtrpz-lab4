package listener

import (
	"encoding/json"
	"github.com/rmerezha/mtrpz-lab4/config"
	"github.com/rmerezha/mtrpz-lab4/runner"
	"log"
	"net/http"
	"sync"
	"time"
)

type PollingListener struct {
	MasterURL string
	Host      string
	Runner    runner.Runner

	mu           sync.Mutex
	Store        *ContainerStateStore
	pollInterval time.Duration

	Token string
}

func NewPollingListener(masterURL, host string, r runner.Runner, interval time.Duration, token string, store *ContainerStateStore) *PollingListener {
	return &PollingListener{
		MasterURL:    masterURL,
		Host:         host,
		Runner:       r,
		Store:        store,
		pollInterval: interval,
		Token:        token,
	}
}

func (pl *PollingListener) Listen(stopCh <-chan struct{}) {
	ticker := time.NewTicker(pl.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			pl.checkAndApply()
		}
	}
}

func (pl *PollingListener) checkAndApply() {
	url := pl.MasterURL + "/api/v1/container?host=" + pl.Host

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("PollingListener: failed to create request: %v", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+pl.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("PollingListener: failed to GET containers: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("PollingListener: unexpected status code %d", resp.StatusCode)
		return
	}

	var containers []config.ContainerStatus
	if err := json.NewDecoder(resp.Body).Decode(&containers); err != nil {
		log.Printf("PollingListener: failed to decode response: %v", err)
		return
	}

	pl.mu.Lock()
	defer pl.mu.Unlock()

	for _, cs := range containers {
		prevState, known := pl.store.Get(cs.Config.Name)
		if !known || prevState != cs.State {
			log.Printf("PollingListener: container %s state changed from %s to %s", cs.Config.Name, prevState, cs.State)
			pl.store.Set(cs.Config.Name, cs.State)

			pl.applyState(cs)
		}
	}
}

func (pl *PollingListener) applyState(cs config.ContainerStatus) {
	name := cs.Config.Name

	switch cs.State {
	case config.StateNew:
		if err := pl.Runner.PullImage(cs.Config.Image); err != nil {
			log.Printf("Runner.PullImage error for %s: %v", name, err)
		}
		if err := pl.Runner.Run(cs.Config); err != nil {
			log.Printf("Runner.Run error for %s: %v", name, err)
		}
		cs.State = config.StateRunning
	case config.StatePaused:
		// TODO
		log.Println("not implemented yet")
	case config.StateRestarting:
		if err := pl.Runner.Restart(name); err != nil {
			log.Printf("Runner.Restart error for %s: %v", name, err)
		}
	case config.StateRemoving:
		if err := pl.Runner.Remove(name); err != nil {
			log.Printf("Runner.Remove error for %s: %v", name, err)
		}
	case config.StateExited:
		if err := pl.Runner.Stop(name); err != nil {
			log.Printf("Runner.Stop error for %s: %v", name, err)
		}
	case config.StateDead:
		if err := pl.Runner.Kill(name); err != nil {
			log.Printf("Runner.Kill error for %s: %v", name, err)
		}
	default:
		log.Printf("PollingListener: unknown state %s for container %s", cs.State, name)
	}
}
