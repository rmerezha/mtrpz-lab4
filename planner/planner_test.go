package planner

import (
	"testing"

	"github.com/rmerezha/mtrpz-lab4/config"
)

func setupPlanner() *Planner {
	manifest := &config.Manifest{
		Name: "example",
		Containers: []config.Container{
			{Name: "web", Host: "node1", Image: "nginx"},
			{Name: "app", Host: "node1", Image: "myapp"},
			{Name: "db", Host: "node2", Image: "postgres"},
		},
	}
	return NewPlanner(manifest)
}

func TestUpdateState_Success(t *testing.T) {
	p := setupPlanner()

	ok := p.UpdateState("node1", "web", StateRunning)
	if !ok {
		t.Fatal("expected UpdateState to return true")
	}

	containers := p.ListContainersByHost("node1")
	found := false
	for _, c := range containers {
		if c.Config.Name == "web" {
			found = true
			if c.State != StateRunning {
				t.Errorf("expected state to be %q, got %q", StateRunning, c.State)
			}
		}
	}
	if !found {
		t.Error("container 'web' not found on node1")
	}
}

func TestUpdateState_FailWrongHost(t *testing.T) {
	p := setupPlanner()

	ok := p.UpdateState("node3", "web", StateRunning)
	if ok {
		t.Error("expected UpdateState to return false for unknown host")
	}
}

func TestUpdateState_FailWrongContainer(t *testing.T) {
	p := setupPlanner()

	ok := p.UpdateState("node1", "unknown", StateRunning)
	if ok {
		t.Error("expected UpdateState to return false for unknown container")
	}
}

func TestListContainersByHost(t *testing.T) {
	p := setupPlanner()

	list := p.ListContainersByHost("node1")
	if len(list) != 2 {
		t.Errorf("expected 2 containers on node1, got %d", len(list))
	}

	names := map[string]bool{}
	for _, c := range list {
		names[c.Config.Name] = true
	}

	if !names["web"] || !names["app"] {
		t.Errorf("expected containers 'web' and 'app' on node1, got %+v", names)
	}
}

func TestListContainersByHost_Empty(t *testing.T) {
	p := setupPlanner()

	list := p.ListContainersByHost("unknown-host")
	if len(list) != 0 {
		t.Errorf("expected 0 containers, got %d", len(list))
	}
}
