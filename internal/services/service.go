// Package services contains the service domain model and persistence logic.
package services

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Service represents a deployable workload managed by Platform-Infra.
type Service struct {
	ID        uuid.UUID `json:"id"`
	ProjectID uuid.UUID `json:"project_id"`
	Name      string    `json:"name"`
	Image     string    `json:"image"`
	Replicas  int       `json:"replicas"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository defines persistence operations for services.
type Repository interface {
	Create(ctx context.Context, service *Service) error
	Get(ctx context.Context, id uuid.UUID) (*Service, error)
	ListByProject(ctx context.Context, projectID uuid.UUID) ([]Service, error)
}
