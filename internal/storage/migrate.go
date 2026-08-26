package storage

import (
	"context"
	"fmt"
)

// Migration contains a single database migration.
type Migration struct {
	Version string
	SQL     string
}

// Migrate applies all pending database migrations.
func (p *Postgres) Migrate(ctx context.Context) error {
	if _, err := p.Pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}

	migrations := []Migration{
		{
			Version: "001_initial",
			SQL: `
				CREATE TABLE IF NOT EXISTS projects (
					id UUID PRIMARY KEY,
					name TEXT NOT NULL UNIQUE,
					description TEXT NOT NULL DEFAULT '',
					created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
				);

				CREATE TABLE IF NOT EXISTS services (
					id UUID PRIMARY KEY,
					project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
					name TEXT NOT NULL,
					image TEXT NOT NULL,
					replicas INTEGER NOT NULL DEFAULT 1,
					created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

					UNIQUE(project_id, name)
				);

				CREATE TABLE IF NOT EXISTS deployments (
					id UUID PRIMARY KEY,
					service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
					version TEXT NOT NULL,
					status TEXT NOT NULL DEFAULT 'pending',
					created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
					updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
				);

				CREATE INDEX IF NOT EXISTS idx_services_project_id
					ON services(project_id);

				CREATE INDEX IF NOT EXISTS idx_deployments_service_id
					ON deployments(service_id);
			`,
		},
	}

	for _, migration := range migrations {
		if err := p.applyMigration(ctx, migration); err != nil {
			return err
		}
	}

	return nil
}

func (p *Postgres) applyMigration(
	ctx context.Context,
	migration Migration,
) error {
	var exists bool

	err := p.Pool.QueryRow(
		ctx,
		`
		SELECT EXISTS(
			SELECT 1
			FROM schema_migrations
			WHERE version = $1
		)
		`,
		migration.Version,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf(
			"check migration %s: %w",
			migration.Version,
			err,
		)
	}

	if exists {
		return nil
	}

	tx, err := p.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf(
			"begin migration %s: %w",
			migration.Version,
			err,
		)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if _, err := tx.Exec(ctx, migration.SQL); err != nil {
		return fmt.Errorf(
			"execute migration %s: %w",
			migration.Version,
			err,
		)
	}

	if _, err := tx.Exec(
		ctx,
		`
		INSERT INTO schema_migrations (version)
		VALUES ($1)
		`,
		migration.Version,
	); err != nil {
		return fmt.Errorf(
			"record migration %s: %w",
			migration.Version,
			err,
		)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf(
			"commit migration %s: %w",
			migration.Version,
			err,
		)
	}

	return nil
}
