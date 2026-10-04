package database

import (
	"context"
	"fmt"
	"time"

	"backend-alusi-go/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// PostgresDB wraps pgxpool.Pool
type PostgresDB struct {
	Pool *pgxpool.Pool
}

// NewPostgresDB establishes a connection pool to PostgreSQL with retries
func NewPostgresDB(ctx context.Context, cfg *config.DatabaseConfig) (*PostgresDB, error) {
	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse postgres config: %w", err)
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MinConns = int32(cfg.MaxIdleConns)
	poolConfig.MaxConnIdleTime = cfg.MaxIdleTime
	poolConfig.MaxConnLifetime = 1 * time.Hour
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	var pool *pgxpool.Pool
	maxRetries := 5
	for attempt := 1; attempt <= maxRetries; attempt++ {
		log.Info().
			Str("host", cfg.Host).
			Str("port", cfg.Port).
			Str("database", cfg.Name).
			Int("attempt", attempt).
			Msg("Connecting to PostgreSQL database...")

		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			err = pool.Ping(pingCtx)
			cancel()
			if err == nil {
				log.Info().Msg("Successfully connected and pinged PostgreSQL database")
				return &PostgresDB{Pool: pool}, nil
			}
			pool.Close()
		}

		log.Warn().
			Err(err).
			Int("attempt", attempt).
			Msg("PostgreSQL not ready yet, retrying in 2 seconds...")

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
		}
	}

	return nil, fmt.Errorf("could not connect to PostgreSQL after %d attempts: %w", maxRetries, err)
}

// Close gracefully closes the database pool
func (db *PostgresDB) Close() {
	if db != nil && db.Pool != nil {
		log.Info().Msg("Closing PostgreSQL connection pool")
		db.Pool.Close()
	}
}

// Ping checks if database connection is alive
func (db *PostgresDB) Ping(ctx context.Context) error {
	if db == nil || db.Pool == nil {
		return fmt.Errorf("database connection is nil")
	}
	return db.Pool.Ping(ctx)
}
