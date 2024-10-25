package postgres

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"

	"github.com/payment-api/infrastructure/exceptions"
	"github.com/payment-api/infrastructure/logger"
)

type Repository struct {
	DB *sql.DB
}

func (r *Repository) GetById(query string, id string, args ...interface{}) error {
	row := r.DB.QueryRow(query, id)

	err := row.Scan(args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Push(query string, args ...interface{}) error {
	row, err := r.DB.Exec(query, args...)
	if err != nil {
		return err
	}

	count, err := row.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return exceptions.EntityNotFoundError
	}

	return nil
}

func connectPostgresDB(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Error(logger.ConfigError, err)

		return nil, err
	}

	return db, pingSql(ctx, db)
}

func pingSql(ctx context.Context, db *sql.DB) (err error) {
	for start := time.Now(); time.Since(start) < (5 * time.Second); {
		err = db.PingContext(ctx)
		if err == nil {
			break
		}

		time.Sleep(1 * time.Second)
	}
	return err
}

func NewRepository(ctx context.Context, dsn string) (*Repository, error) {
	db, err := connectPostgresDB(ctx, dsn)
	if err != nil {
		return nil, err
	}

	return &Repository{
		DB: db,
	}, nil
}
