package services

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"github.com/ikwukao/platform-infra/internal/storage"
)

// PostgresRepository implements Repository using PostgreSQL.
type PostgresRepository struct {
	db *storage.Postgres
}

// NewPostgresRepository creates a service repository backed by PostgreSQL.
func NewPostgresRepository(db *storage.Postgres) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Create persists a new service.
func (r *PostgresRepository) Create(
	ctx context.Context,
	service *Service,
) error {
	_, err := r.db.Pool.Exec(
		ctx,
		`
		INSERT INTO services (
			id,
			project_id,
			name,
			image,
			replicas
		)
		VALUES ($1, $2, $3, $4, $5)
		`,
		service.ID,
		service.ProjectID,
		service.Name,
		service.Image,
		service.Replicas,
	)

	if err != nil {
		return fmt.Errorf("create service: %w", err)
	}

	return nil
}

// Get retrieves a service by ID.
func (r *PostgresRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*Service, error) {
	service := &Service{}

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			project_id,
			name,
			image,
			replicas,
			created_at,
			updated_at
		FROM services
		WHERE id = $1
		`,
		id,
	).Scan(
		&service.ID,
		&service.ProjectID,
		&service.Name,
		&service.Image,
		&service.Replicas,
		&service.CreatedAt,
		&service.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("get service: %w", err)
	}

	return service, nil
}

// ListByProject returns all services belonging to a project.
func (r *PostgresRepository) ListByProject(
	ctx context.Context,
	projectID uuid.UUID,
) ([]Service, error) {
	rows, err := r.db.Pool.Query(
		ctx,
		`
		SELECT
			id,
			project_id,
			name,
			image,
			replicas,
			created_at,
			updated_at
		FROM services
		WHERE project_id = $1
		ORDER BY created_at DESC
		`,
		projectID,
	)
	if err != nil {
		return nil, fmt.Errorf("list services: %w", err)
	}
	defer rows.Close()

	var services []Service

	for rows.Next() {
		var service Service

		if err := rows.Scan(
			&service.ID,
			&service.ProjectID,
			&service.Name,
			&service.Image,
			&service.Replicas,
			&service.CreatedAt,
			&service.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan service: %w", err)
		}

		services = append(services, service)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate services: %w", err)
	}

	return services, nil
}
