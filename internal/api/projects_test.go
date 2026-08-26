package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/ikwukao/platform-infra/internal/projects"
)

type mockProjectRepository struct {
	items []projects.Project
}

func (m *mockProjectRepository) Create(
	_ context.Context,
	project *projects.Project,
) error {
	m.items = append(m.items, *project)

	return nil
}

func (m *mockProjectRepository) Get(
	_ context.Context,
	id uuid.UUID,
) (*projects.Project, error) {
	for _, project := range m.items {
		if project.ID == id {
			result := project
			return &result, nil
		}
	}

	return nil, projects.ErrNotFound
}

func (m *mockProjectRepository) List(
	_ context.Context,
) ([]projects.Project, error) {
	return m.items, nil
}

func TestProjectHandlerCreate(t *testing.T) {
	repository := &mockProjectRepository{}
	handler := NewProjectHandler(repository)

	requestBody := `{
		"name": "flux-platform",
		"description": "Platform infrastructure"
	}`

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/projects",
		strings.NewReader(requestBody),
	)

	request.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	handler.create(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}

	if len(repository.items) != 1 {
		t.Fatalf(
			"expected 1 project, got %d",
			len(repository.items),
		)
	}

	project := repository.items[0]

	if project.Name != "flux-platform" {
		t.Fatalf(
			"expected project name %q, got %q",
			"flux-platform",
			project.Name,
		)
	}

	if project.Description != "Platform infrastructure" {
		t.Fatalf(
			"expected description %q, got %q",
			"Platform infrastructure",
			project.Description,
		)
	}
}

func TestProjectHandlerCreateRequiresName(t *testing.T) {
	repository := &mockProjectRepository{}
	handler := NewProjectHandler(repository)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/projects",
		strings.NewReader(`{
			"description": "Missing project name"
		}`),
	)

	recorder := httptest.NewRecorder()

	handler.create(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusBadRequest,
			recorder.Code,
		)
	}

	if len(repository.items) != 0 {
		t.Fatalf("project should not have been created")
	}
}

func TestProjectHandlerGet(t *testing.T) {
	projectID := uuid.New()

	repository := &mockProjectRepository{
		items: []projects.Project{
			{
				ID:          projectID,
				Name:        "flux-platform",
				Description: "Platform infrastructure",
			},
		},
	}

	handler := NewProjectHandler(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/projects/"+projectID.String(),
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.get(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}

func TestProjectHandlerGetNotFound(t *testing.T) {
	repository := &mockProjectRepository{}
	handler := NewProjectHandler(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/projects/"+uuid.New().String(),
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.get(recorder, request)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusNotFound,
			recorder.Code,
		)
	}
}

func TestProjectHandlerList(t *testing.T) {
	repository := &mockProjectRepository{
		items: []projects.Project{
			{
				ID:   uuid.New(),
				Name: "flux-platform",
			},
			{
				ID:   uuid.New(),
				Name: "strata-platform",
			},
		},
	}

	handler := NewProjectHandler(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/projects",
		nil,
	)

	recorder := httptest.NewRecorder()

	handler.list(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}

func TestServerProjectRoutes(t *testing.T) {
	repository := &mockProjectRepository{}
	server := NewServer(repository)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/projects",
		strings.NewReader(`{
			"name": "platform-infra"
		}`),
	)

	recorder := httptest.NewRecorder()

	server.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}
}

func TestHealthEndpoint(t *testing.T) {
	repository := &mockProjectRepository{}
	server := NewServer(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/healthz",
		nil,
	)

	recorder := httptest.NewRecorder()

	server.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}

func TestReadyEndpoint(t *testing.T) {
	repository := &mockProjectRepository{}
	server := NewServer(repository)

	request := httptest.NewRequest(
		http.MethodGet,
		"/readyz",
		nil,
	)

	recorder := httptest.NewRecorder()

	server.Handler().ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusOK,
			recorder.Code,
		)
	}
}
