// Package storage provides persistence primitives for Platform-Infra.
package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Postgres manages the PostgreSQL connection pool used by Platform-Infra.
type Postgres struct {
	Pool *pgxpool.Pool
}

// NewPostgres creates a PostgreSQL connection pool using the supplied
// connection string.
func NewPostgres(ctx context.Context, databaseURL string) (*Postgres, error) {
	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &Postgres{
		Pool: pool,
	}, nil
}

// Close releases all connections managed by the PostgreSQL pool.
func (p *Postgres) Close() {
	p.Pool.Close()
}
