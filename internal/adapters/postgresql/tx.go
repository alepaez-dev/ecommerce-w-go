package postgresql

import (
	"context"

	repo "github.com/alepaez-dev/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

type TxManager struct {
	db *pgx.Conn
}

func NewTxManager(db *pgx.Conn) *TxManager {
	return &TxManager{db: db}
}

func (m *TxManager) WithTx(ctx context.Context, fn func(repo.Querier) error) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := repo.New(tx)
	if err := fn(qtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
