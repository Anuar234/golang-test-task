package storage

import (
	"context"
	"database/sql"
)

type NumbersStore interface {
	// добавляет число и сразу возвращает сортированный список
	AddAndList(ctx context.Context, value int64) ([]int64, error)
	// гарантирует, что таблица есть
	EnsureSchema(ctx context.Context) error
}

type PostgresStore struct {
	// держим ссылку на db пул
	db *sql.DB
}

func NewPostgres(db *sql.DB) *PostgresStore {
	// простой конструктор без магии
	return &PostgresStore{db: db}
}

func (s *PostgresStore) EnsureSchema(ctx context.Context) error {
	// минимальная миграция, чтобы не тащить внешние тулзы
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS numbers (
			id BIGSERIAL PRIMARY KEY,
			value BIGINT NOT NULL
		);
	`)
	return err
}

func (s *PostgresStore) AddAndList(ctx context.Context, value int64) ([]int64, error) {
	// все в одной транзакции, чтобы список точно включал вставку
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// пишем новое число
	if _, err := tx.ExecContext(ctx, `INSERT INTO numbers (value) VALUES ($1)`, value); err != nil {
		return nil, err
	}

	// читаем весь список отсортированным
	rows, err := tx.QueryContext(ctx, `SELECT value FROM numbers ORDER BY value ASC, id ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// собираем значения в слайс
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

	// фиксируем транзакцию
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return values, nil
}
