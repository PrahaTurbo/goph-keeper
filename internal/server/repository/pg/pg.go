// Package pg provides functions for creating and managing PostgreSQL connections.
package pg

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/PrahaTurbo/goph-keeper/internal/server/config"
)

// DefaultQueryTimeout specifies a predefined period of time to wait for a query to complete before cancelling it.
const DefaultQueryTimeout = time.Second * 5

// NewPGPool creates a new PostgreSQL connection pool using the given configuration.
func NewPGPool(cfg config.PG) (*pgxpool.Pool, error) {
	poolConfig, err := newPGPoolConfig(cfg)
	if err != nil {
		return nil, err
	}

	conn, err := newPGConnection(poolConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func newPGPoolConfig(cfg config.PG) (*pgxpool.Config, error) {
	dsn := fmt.Sprintf(
		"%s://%s:%s@%s:%d/%s?sslmode=prefer",
		"postgres",
		url.QueryEscape(cfg.User),
		url.QueryEscape(cfg.Password),
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	return poolConfig, nil
}

func newPGConnection(poolConfig *pgxpool.Config) (*pgxpool.Pool, error) {
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
