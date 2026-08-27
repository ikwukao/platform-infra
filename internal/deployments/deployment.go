// Package deployments contains the deployment domain model and persistence logic.
package deployments

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Deployment status values represent the lifecycle of a deployment.
const (
	StatusPending   = "pending"
	StatusRunning   = "running"
	StatusSucceeded = "succeeded"
	StatusFailed    = "failed"
)

// Deployment represents a rollout of a service managed by Platform-Infra.
type Deployment struct {
	ID        uuid.UUID `json:"id"`
	ServiceID uuid.UUID `json:"service_id"`
	Version   string    `json:"version"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Repository defines persistence operations for deployments.
type Repository interface {
	Create(ctx context.Context, deployment *Deployment) error
	Get(ctx context.Context, id uuid.UUID) (*Deployment, error)
	ListByService(ctx context.Context, serviceID uuid.UUID) ([]Deployment, error)
	ListPending(ctx context.Context) ([]Deployment, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}
