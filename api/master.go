package api

import (
	"encoding/json"
	"github.com/rmerezha/mtrpz-lab4/auth"
	"net/http"

	"github.com/rmerezha/mtrpz-lab4/planner"
)

type Server struct {
	Planner *planner.Planner
	Auth    *auth.Manager
}

type StateUpdateRequest struct {
	Host          string `json:"host"`
	ContainerName string `json:"name"`
	State         string `json:"state"`
}

func (s *Server) handleUpdateState(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req StateUpdateRequest
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

func (s *Server) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/state", withAuth(s.Auth, s.handleUpdateState))
	mux.HandleFunc("/api/v1/container", withAuth(s.Auth, s.handleListContainers))
}
