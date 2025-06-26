package store

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:generate options-gen -out-filename=storage_options.gen.go -from-struct=Options
type Options struct {
	dbName        string `option:"mandatory"`
	dbUser        string `option:"mandatory"`
	dbPassword    string `option:"mandatory"`
	dbHost        string `option:"mandatory"`
	dbPort        string `option:"mandatory"`
	dbSSLMode     string `option:"optional" default:"disable"`
	dbSSLRootCert string
	dbSSLKey      string
}

// Storage is the database connection pool.
type Storage struct {
	db *pgxpool.Pool
}

// NewStorage creates a new database connection pool.
func NewStorage(ctx context.Context, opts Options) (*Storage, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid options for storage: %v", err)
	}

	dbconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		opts.dbHost, opts.dbPort, opts.dbUser, opts.dbPassword, opts.dbName, opts.dbSSLMode)

	pool, err := pgxpool.New(ctx, dbconn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &Storage{db: pool}, nil
}

// Close closes the database connection pool.
func (s *Storage) Close() {
	s.db.Close()
}

// Ping the database
func (s *Storage) Ping(ctx context.Context) error {
	return s.db.Ping(ctx)
}
