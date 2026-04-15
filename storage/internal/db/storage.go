package db

import (
	"context"
	"fmt"
	"techsupport/storage/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	Pool *pgxpool.Pool
}

func NewStorage(cfg config.StorageConfig) (*Storage, error) {
    dsn := config.GetConfig().BuildDSN()

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    poolConfig, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        return nil, fmt.Errorf("failed to parse dsn: %w", err)
    }

    poolConfig.MaxConns = cfg.MaxConns
    poolConfig.MinConns = cfg.MinConns
    poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime
    poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
    poolConfig.HealthCheckPeriod = cfg.HealthCheck

    pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
    if err != nil {
        return nil, fmt.Errorf("failed to create pool: %w", err)
    }

    if err := pool.Ping(ctx); err != nil {
        return nil, fmt.Errorf("db ping failed: %w", err)
    }

    return &Storage{Pool: pool}, nil
}