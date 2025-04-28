package postgres

import (
	"context"

	"github.com/cektrendstudio/cektrend-engine-go/service"
	"github.com/jmoiron/sqlx"
)

type sqlRepository struct {
	DB *sqlx.DB
}

func NewSQLRepository(db *sqlx.DB) service.SQLRepository {
	return sqlRepository{DB: db}
}

func (repo sqlRepository) BeginTxx() (tx *sqlx.Tx, err error) {
	return repo.DB.BeginTxx(context.Background(), nil)
}

func (repo sqlRepository) Commit(tx *sqlx.Tx) (err error) {
	return tx.Commit()
}

func (repo sqlRepository) Rollback(tx *sqlx.Tx) (err error) {
	return tx.Rollback()
}
