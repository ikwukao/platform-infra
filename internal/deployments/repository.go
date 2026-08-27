package deployments

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

// NewPostgresRepository creates a deployment repository backed by PostgreSQL.
func NewPostgresRepository(db *storage.Postgres) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Create persists a new deployment.
func (r *PostgresRepository) Create(
	ctx context.Context,
	deployment *Deployment,
) error {
	_, err := r.db.Pool.Exec(
		ctx,
		`
		INSERT INTO deployments (
			id,
			service_id,
			version,
			status
		)
		VALUES ($1, $2, $3, $4)
		`,
		deployment.ID,
		deployment.ServiceID,
		deployment.Version,
		deployment.Status,
	)
	if err != nil {
		return fmt.Errorf("create deployment: %w", err)
	}

	return nil
}

// Get retrieves a deployment by ID.
func (r *PostgresRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*Deployment, error) {
	deployment := &Deployment{}

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			service_id,
			version,
			status,
			created_at,
			updated_at
		FROM deployments
		WHERE id = $1
		`,
		id,
	).Scan(
		&deployment.ID,
		&deployment.ServiceID,
		&deployment.Version,
		&deployment.Status,
		&deployment.CreatedAt,
		&deployment.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("get deployment: %w", err)
	}

	return deployment, nil
}

// ListByService returns all deployments belonging to a service.
func (r *PostgresRepository) ListByService(
	ctx context.Context,
	serviceID uuid.UUID,
) ([]Deployment, error) {
	rows, err := r.db.Pool.Query(
		ctx,
		`
		SELECT
			id,
			service_id,
			version,
			status,
			created_at,
			updated_at
		FROM deployments
		WHERE service_id = $1
		ORDER BY created_at DESC
		`,
		serviceID,
	)
	if err != nil {
		return nil, fmt.Errorf("list deployments: %w", err)
	}
	defer rows.Close()

	var deployments []Deployment

	for rows.Next() {
		var deployment Deployment

		if err := rows.Scan(
			&deployment.ID,
			&deployment.ServiceID,
			&deployment.Version,
			&deployment.Status,
			&deployment.CreatedAt,
			&deployment.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan deployment: %w", err)
		}

		deployments = append(deployments, deployment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate deployments: %w", err)
	}

	return deployments, nil
}
