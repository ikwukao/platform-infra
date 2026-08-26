// Package projects contains the project domain model and persistence logic.
package projects

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Project represents an application or infrastructure project managed by
// Platform-Infra.
type Project struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Repository defines persistence operations for projects.
type Repository interface {
	Create(ctx context.Context, project *Project) error
	Get(ctx context.Context, id uuid.UUID) (*Project, error)
	List(ctx context.Context) ([]Project, error)
}
