package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/ikwukao/platform-infra/internal/deployments"
)

// DeploymentHandler handles deployment API operations.
type DeploymentHandler struct {
	repository deployments.Repository
}

// NewDeploymentHandler creates a deployment HTTP handler.
func NewDeploymentHandler(
	repository deployments.Repository,
) *DeploymentHandler {
	return &DeploymentHandler{
		repository: repository,
	}
}

type createDeploymentRequest struct {
	Version string `json:"version"`
	Status  string `json:"status"`
}

func (h *DeploymentHandler) create(
	w http.ResponseWriter,
	r *http.Request,
) {
	serviceIDValue := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/services/",
	)

	serviceIDValue = strings.TrimSuffix(
		serviceIDValue,
		"/deployments",
	)

	serviceID, err := uuid.Parse(serviceIDValue)
	if err != nil {
		http.Error(
			w,
			"invalid service id",
			http.StatusBadRequest,
		)
		return
	}

	var request createDeploymentRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}

	request.Version = strings.TrimSpace(request.Version)
	request.Status = strings.TrimSpace(request.Status)

	if request.Version == "" {
		http.Error(
			w,
			"deployment version is required",
			http.StatusBadRequest,
		)
		return
	}

	if request.Status == "" {
		request.Status = "pending"
	}

	deployment := &deployments.Deployment{
		ID:        uuid.New(),
		ServiceID: serviceID,
		Version:   request.Version,
		Status:    request.Status,
	}

	if err := h.repository.Create(
		r.Context(),
		deployment,
	); err != nil {
		http.Error(
			w,
			"failed to create deployment",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusCreated, deployment)
}

func (h *DeploymentHandler) get(
	w http.ResponseWriter,
	r *http.Request,
) {
	idValue := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/deployments/",
	)

	id, err := uuid.Parse(idValue)
	if err != nil {
		http.Error(
			w,
			"invalid deployment id",
			http.StatusBadRequest,
		)
		return
	}

	deployment, err := h.repository.Get(
		r.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, deployments.ErrNotFound) {
			http.Error(
				w,
				"deployment not found",
				http.StatusNotFound,
			)
			return
		}

		http.Error(
			w,
			"failed to retrieve deployment",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, deployment)
}

func (h *DeploymentHandler) listByService(
	w http.ResponseWriter,
	r *http.Request,
) {
	serviceIDValue := strings.TrimPrefix(
		r.URL.Path,
		"/api/v1/services/",
	)

	serviceIDValue = strings.TrimSuffix(
		serviceIDValue,
		"/deployments",
	)

	serviceID, err := uuid.Parse(serviceIDValue)
	if err != nil {
		http.Error(
			w,
			"invalid service id",
			http.StatusBadRequest,
		)
		return
	}

	items, err := h.repository.ListByService(
		r.Context(),
		serviceID,
	)
	if err != nil {
		http.Error(
			w,
			"failed to list deployments",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, items)
}
