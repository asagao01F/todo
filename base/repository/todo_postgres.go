package repository

import (
	"context"
	"database/sql"
	"todo/base/model"
)

// インターフェースの指定（implements等）は一切書かない
type PostgresTodoRepository struct {
	db *sql.DB
}

func NewPostgresTodoRepository(db *sql.DB) *PostgresTodoRepository {
	return &PostgresTodoRepository{db: db}
}

func (r *PostgresTodoRepository) Create(ctx context.Context, todo *model.Todo) error {
	// SQL処理...
	return nil
}
