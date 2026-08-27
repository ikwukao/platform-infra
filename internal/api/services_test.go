package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/ikwukao/platform-infra/internal/services"
)

type mockServiceRepository struct {
	items []services.Service
}

func (m *mockServiceRepository) Create(
	_ context.Context,
	service *services.Service,
) error {
	m.items = append(m.items, *service)

	return nil
}

func (m *mockServiceRepository) Get(
	_ context.Context,
	id uuid.UUID,
) (*services.Service, error) {
	for _, service := range m.items {
		if service.ID == id {
			result := service
			return &result, nil
		}
	}

	return nil, services.ErrNotFound
}

func (m *mockServiceRepository) ListByProject(
	_ context.Context,
	projectID uuid.UUID,
) ([]services.Service, error) {
	var result []services.Service

	for _, service := range m.items {
		if service.ProjectID == projectID {
			result = append(result, service)
		}
	}

	return result, nil
}

func TestServiceHandlerCreate(t *testing.T) {
	projectID := uuid.New()

	repository := &mockServiceRepository{}
	handler := NewServiceHandler(repository)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/projects/"+projectID.String()+"/services",
		strings.NewReader(`{
			"name": "gateway",
			"image": "ghcr.io/ikwukao/flux-gateway:latest",
			"replicas": 2
		}`),
	)

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
		t.Fatalf("expected 1 service, got %d", len(repository.items))
	}

	service := repository.items[0]

	if service.Name != "gateway" {
		t.Fatalf(
			"expected service name %q, got %q",
			"gateway",
			service.Name,
		)
	}

	if service.Replicas != 2 {
		t.Fatalf(
			"expected 2 replicas, got %d",
			service.Replicas,
		)
	}
}

func TestServiceHandlerDefaultsReplicas(t *testing.T) {
	projectID := uuid.New()

	repository := &mockServiceRepository{}
	handler := NewServiceHandler(repository)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/projects/"+projectID.String()+"/services",
		strings.NewReader(`{
			"name": "gateway",
			"image": "ghcr.io/ikwukao/flux-gateway:latest"
		}`),
	)

	recorder := httptest.NewRecorder()

	handler.create(recorder, request)

	if recorder.Code != http.StatusCreated {
		t.Fatalf(
			"expected status %d, got %d",
			http.StatusCreated,
			recorder.Code,
		)
	}

	if repository.items[0].Replicas != 1 {
		t.Fatalf(
			"expected default replicas 1, got %d",
			repository.items[0].Replicas,
		)
	}
}

func TestServiceHandlerRequiresImage(t *testing.T) {
	projectID := uuid.New()

	repository := &mockServiceRepository{}
	handler := NewServiceHandler(repository)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/projects/"+projectID.String()+"/services",
		strings.NewReader(`{
			"name": "gateway"
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
}
