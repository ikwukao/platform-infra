// Package api provides the HTTP API for the Platform-Infra controller.
package api

import (
	"encoding/json"
	"net/http"
)

// Server exposes the Platform-Infra HTTP API.
type Server struct {
	mux *http.ServeMux
}

// NewServer creates a new API server with the default platform routes.
func NewServer() *Server {
	server := &Server{
		mux: http.NewServeMux(),
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
