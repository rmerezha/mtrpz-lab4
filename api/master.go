package api

import (
	"encoding/json"
	"github.com/rmerezha/mtrpz-lab4/auth"
	"github.com/rmerezha/mtrpz-lab4/config"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"

	"github.com/rmerezha/mtrpz-lab4/planner"
)

type Server struct {
	Planner *planner.Planner
	Auth    *auth.Manager
}

func (s *Server) handleUpdateState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Host          string `json:"host"`
		ContainerName string `json:"name"`
		State         string `json:"state"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	ok := s.Planner.UpdateState(req.Host, req.ContainerName, req.State)
	if !ok {
		http.Error(w, "container not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) handleListContainers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	host := r.URL.Query().Get("host")
	if host == "" {
		http.Error(w, "missing 'host' query param", http.StatusBadRequest)
		return
	}

	containers := s.Planner.ListContainersByHost(host)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(containers); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func (s *Server) handleContainerAction(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Action    string `json:"action"`
		Host      string `json:"host"`
		Container string `json:"container"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	var targetState planner.ContainerState

	switch req.Action {
	case "stop":
		targetState = planner.StateExited
	case "kill":
		targetState = planner.StateDead
	case "restart":
		targetState = planner.StateRestarting
	case "rm":
		targetState = planner.StateRemoving
	default:
		http.Error(w, "unsupported action", http.StatusBadRequest)
		return
	}

	if ok := s.Planner.UpdateState(req.Host, req.Container, targetState); !ok {
		http.Error(w, "container not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleManifestUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var manifest config.Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		http.Error(w, "invalid YAML: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err := manifest.Validate(); err != nil {
		http.Error(w, "invalid manifest: "+err.Error(), http.StatusBadRequest)
		return
	}

	s.Planner.AddManifest(&manifest)
	w.WriteHeader(http.StatusCreated)
}

func (s *Server) handleManifestDown(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Manifest string `json:"manifest"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Manifest == "" {
		http.Error(w, "missing manifest name", http.StatusBadRequest)
		return
	}

	ok := s.Planner.MarkManifestRemoving(req.Manifest)
	if !ok {
		http.Error(w, "manifest not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/state", withAuth(s.Auth, s.handleUpdateState))
	mux.HandleFunc("/api/v1/container", withAuth(s.Auth, s.handleListContainers))
	mux.HandleFunc("/api/v1/container/action", withAuth(s.Auth, s.handleContainerAction))
	mux.HandleFunc("/api/v1/manifest/up", withAuth(s.Auth, s.handleManifestUp))
}
