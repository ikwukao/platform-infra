// Package api provides the HTTP API for the Platform-Infra controller.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/ikwukao/platform-infra/internal/projects"
)

// Server exposes the Platform-Infra HTTP API.
type Server struct {
	mux            *http.ServeMux
	projectHandler *ProjectHandler
}

// NewServer creates a new API server with the default platform routes.
func NewServer(projectRepository projects.Repository) *Server {
	server := &Server{
		mux:            http.NewServeMux(),
		projectHandler: NewProjectHandler(projectRepository),
	}

	server.registerRoutes()

	return server
}

// Handler returns the HTTP handler for the API server.
func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) registerRoutes() {
	s.mux.HandleFunc("/healthz", s.health)
	s.mux.HandleFunc("/readyz", s.ready)

	s.mux.HandleFunc(
		"/api/v1/projects",
		func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				s.projectHandler.list(w, r)
			case http.MethodPost:
				s.projectHandler.create(w, r)
			default:
				http.Error(
					w,
					"method not allowed",
					http.StatusMethodNotAllowed,
				)
			}
		},
	)

	s.mux.HandleFunc(
		"/api/v1/projects/",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.Error(
					w,
					"method not allowed",
					http.StatusMethodNotAllowed,
				)
				return
			}

			s.projectHandler.get(w, r)
		},
	)
}

func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func (s *Server) ready(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status": "ready",
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}
