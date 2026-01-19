package storage

import (
	"context"
	"database/sql"
)

type NumbersStore interface {
	AddAndList(ctx context.Context, value int64) ([]int64, error)
	EnsureSchema(ctx context.Context) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *PostgresStore {
	return &PostgresStore{db: db}
}

func (s *PostgresStore) EnsureSchema(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS numbers (
			id BIGSERIAL PRIMARY KEY,
			value BIGINT NOT NULL
		);
	`)
	return err
}

func (s *PostgresStore) AddAndList(ctx context.Context, value int64) ([]int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, `INSERT INTO numbers (value) VALUES ($1)`, value); err != nil {
		return nil, err
	}

	rows, err := tx.QueryContext(ctx, `SELECT value FROM numbers ORDER BY value ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values []int64
	for rows.Next() {
		var v int64
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		values = append(values, v)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return values, nil
}
