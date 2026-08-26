package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/ikwukao/platform-infra/internal/services"
)

// ServiceHandler handles service API operations.
type ServiceHandler struct {
	repository services.Repository
}

// NewServiceHandler creates a service HTTP handler.
func NewServiceHandler(
	repository services.Repository,
) *ServiceHandler {
	return &ServiceHandler{
		repository: repository,
	}
}

type createServiceRequest struct {
	Name     string `json:"name"`
	Image    string `json:"image"`
	Replicas int    `json:"replicas"`
}

func (h *ServiceHandler) create(
	w http.ResponseWriter,
	r *http.Request,
) {
	projectIDValue := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/projects/",
	)

	projectIDValue = strings.TrimSuffix(
		projectIDValue,
		"/services",
	)

	projectID, err := uuid.Parse(projectIDValue)
	if err != nil {
		http.Error(
			w,
			"invalid project id",
			http.StatusBadRequest,
		)
		return
	}

	var request createServiceRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	request.Name = strings.TrimSpace(request.Name)
	request.Image = strings.TrimSpace(request.Image)

	if request.Name == "" {
		http.Error(
			w,
			"service name is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.Image == "" {
		http.Error(
			w,
			"service image is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.Replicas < 1 {
		request.Replicas = 1
	}

	service := &services.Service{
		ID:        uuid.New(),
		ProjectID: projectID,
		Name:      request.Name,
		Image:     request.Image,
		Replicas:  request.Replicas,
	}

	if err := h.repository.Create(
		r.Context(),
		service,
	); err != nil {
		http.Error(
			w,
			"failed to create service",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusCreated, service)
}

func (h *ServiceHandler) get(
	w http.ResponseWriter,
	r *http.Request,
) {
	idValue := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/services/",
	)

	id, err := uuid.Parse(idValue)
	if err != nil {
		http.Error(
			w,
			"invalid service id",
			http.StatusBadRequest,
		)
		return
	}

	service, err := h.repository.Get(
		r.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			http.Error(
				w,
				"service not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			"failed to retrieve service",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, service)
}

func (h *ServiceHandler) listByProject(
	w http.ResponseWriter,
	r *http.Request,
) {
	projectIDValue := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/projects/",
	)

	projectIDValue = strings.TrimSuffix(
		projectIDValue,
		"/services",
	)

	projectID, err := uuid.Parse(projectIDValue)
	if err != nil {
		http.Error(
			w,
			"invalid project id",
			http.StatusBadRequest,
		)
		return
	}

	items, err := h.repository.ListByProject(
		r.Context(),
		projectID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to list services",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, items)
}
