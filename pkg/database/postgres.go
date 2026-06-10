package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgresPool(databaseURL string) *pgxpool.Pool {
	cfg, err := pgxpool.ParseConfig(databaseURL)

	if err != nil {
		log.Fatalf("unable to parse database URL: %v", err)
	}

	// Max simultaneous connections
	cfg.MaxConns = 25
	// Connections maintain constantly
	cfg.MinConns = 5
	// Lifetime connection
	cfg.MaxConnLifetime = 1 * time.Hour
	// Close connection after 30min of inactivity
	cfg.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("unable to create connection pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("unable to reach database: %v", err)
	}

	fmt.Println("[INFO] Connected to PostgreSQL")
	return pool
}
