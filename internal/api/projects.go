package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/ikwukao/platform-infra/internal/projects"
)

// ProjectHandler handles project API operations.
type ProjectHandler struct {
	repository projects.Repository
}

// NewProjectHandler creates a project HTTP handler.
func NewProjectHandler(repository projects.Repository) *ProjectHandler {
	return &ProjectHandler{
		repository: repository,
	}
}

type createProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *ProjectHandler) create(
	w http.ResponseWriter,
	r *http.Request,
) {
	var request createProjectRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	request.Name = strings.TrimSpace(request.Name)

	if request.Name == "" {
		http.Error(w, "project name is required", http.StatusBadRequest)
		return
	}

	project := &projects.Project{
		ID:          uuid.New(),
		Name:        request.Name,
		Description: request.Description,
	}

	if err := h.repository.Create(r.Context(), project); err != nil {
		http.Error(
			w,
			"failed to create project",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) get(
	w http.ResponseWriter,
	r *http.Request,
) {
	idValue := strings.TrimPrefix(r.URL.Path, "/api/v1/projects/")

	id, err := uuid.Parse(idValue)
	if err != nil {
		http.Error(w, "invalid project id", http.StatusBadRequest)
		return
	}

	project, err := h.repository.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, projects.ErrNotFound) {
			http.Error(w, "project not found", http.StatusNotFound)
			return
		}

		http.Error(
			w,
			"failed to retrieve project",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) list(
	w http.ResponseWriter,
	r *http.Request,
) {
	items, err := h.repository.List(r.Context())
	if err != nil {
		http.Error(
			w,
			"failed to list projects",
			http.StatusInternalServerError,
		)
		return
	}

	writeJSON(w, http.StatusOK, items)
}
