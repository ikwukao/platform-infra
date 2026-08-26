package projects

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

// NewPostgresRepository creates a project repository backed by PostgreSQL.
func NewPostgresRepository(db *storage.Postgres) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Create persists a new project.
func (r *PostgresRepository) Create(
	ctx context.Context,
	project *Project,
) error {
	_, err := r.db.Pool.Exec(
		ctx,
		`
		INSERT INTO projects (
			id,
			name,
			description
		)
		VALUES ($1, $2, $3)
		`,
		project.ID,
		project.Name,
		project.Description,
	)

	if err != nil {
		return fmt.Errorf("create project: %w", err)
	}

	return nil
}

// Get retrieves a project by ID.
func (r *PostgresRepository) Get(
	ctx context.Context,
	id uuid.UUID,
) (*Project, error) {
	project := &Project{}

	err := r.db.Pool.QueryRow(
		ctx,
		`
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM projects
		WHERE id = $1
		`,
		id,
	).Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&project.CreatedAt,
		&project.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("get project: %w", err)
	}

	return project, nil
}

// List returns all projects.
func (r *PostgresRepository) List(
	ctx context.Context,
) ([]Project, error) {
	rows, err := r.db.Pool.Query(
		ctx,
		`
		SELECT
			id,
			name,
			description,
			created_at,
			updated_at
		FROM projects
		ORDER BY created_at DESC
		`,
	)
	if err != nil {
		return nil, fmt.Errorf("list projects: %w", err)
	}
	defer rows.Close()

	var projects []Project

	for rows.Next() {
		var project Project

		if err := rows.Scan(
			&project.ID,
			&project.Name,
			&project.Description,
			&project.CreatedAt,
			&project.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan project: %w", err)
		}

		projects = append(projects, project)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate projects: %w", err)
	}

	return projects, nil
}
