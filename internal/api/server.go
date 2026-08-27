// Package api provides the HTTP API for Platform-Infra.
package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ikwukao/platform-infra/internal/deployments"
	"github.com/ikwukao/platform-infra/internal/projects"
	"github.com/ikwukao/platform-infra/internal/services"
)

// Server exposes the Platform-Infra HTTP API.
type Server struct {
	mux               *http.ServeMux
	projectHandler    *ProjectHandler
	serviceHandler    *ServiceHandler
	deploymentHandler *DeploymentHandler
}

// NewServer creates an API server backed by the supplied repositories.
func NewServer(
	projectRepository projects.Repository,
	serviceRepository services.Repository,
	deploymentRepository deployments.Repository,
) *Server {
	server := &Server{
		mux:               http.NewServeMux(),
		projectHandler:    NewProjectHandler(projectRepository),
		serviceHandler:    NewServiceHandler(serviceRepository),
		deploymentHandler: NewDeploymentHandler(deploymentRepository),
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
			if strings.HasSuffix(r.URL.Path, "/services") {
				switch r.Method {
				case http.MethodGet:
					s.serviceHandler.listByProject(w, r)

				case http.MethodPost:
					s.serviceHandler.create(w, r)

				default:
					http.Error(
						w,
						"method not allowed",
						http.StatusMethodNotAllowed,
					)
				}

				return
			}

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

	s.mux.HandleFunc(
		"/api/v1/services/",
		func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/deployments") {
				switch r.Method {
				case http.MethodGet:
					s.deploymentHandler.listByService(w, r)

				case http.MethodPost:
					s.deploymentHandler.create(w, r)

				default:
					http.Error(
						w,
						"method not allowed",
						http.StatusMethodNotAllowed,
					)
				}

				return
			}

			if r.Method != http.MethodGet {
				http.Error(
					w,
					"method not allowed",
					http.StatusMethodNotAllowed,
				)
				return
			}

			s.serviceHandler.get(w, r)
		},
	)

	s.mux.HandleFunc(
		"/api/v1/deployments/",
		func(w http.ResponseWriter, r *http.Request) {
			if strings.HasSuffix(r.URL.Path, "/status") {
				if r.Method != http.MethodPatch {
					http.Error(
						w,
						"method not allowed",
						http.StatusMethodNotAllowed,
					)
					return
				}

				s.deploymentHandler.updateStatus(w, r)
				return
			}

			if r.Method != http.MethodGet {
				http.Error(
					w,
					"method not allowed",
					http.StatusMethodNotAllowed,
				)
				return
			}

			s.deploymentHandler.get(w, r)
		},
	)
}

func (s *Server) health(
	w http.ResponseWriter,
	_ *http.Request,
) {
	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "ok",
		},
	)
}

func (s *Server) ready(
	w http.ResponseWriter,
	_ *http.Request,
) {
	writeJSON(
		w,
		http.StatusOK,
		map[string]string{
			"status": "ready",
		},
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	value any,
) {
	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(value)
}
