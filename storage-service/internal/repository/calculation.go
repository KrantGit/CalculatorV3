package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type CalculationRepository interface {
	Save(ctx context.Context, expr string, result float64) error
}

type PostgresRepo struct {
	db *sql.DB
}

func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Save(ctx context.Context, expr string, result float64) error {
	const query = `INSERT INTO calculations (expression, result) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, expr, result)
	if err != nil {
		return fmt.Errorf("postgres save failed: %w", err)
	}
	return nil
}
