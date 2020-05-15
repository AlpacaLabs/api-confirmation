package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

type Client interface {
	RunInTransaction(ctx context.Context, fn func(context.Context, Transaction) error) error
}

type clientImpl struct {
	db *pgx.Conn
}

func NewClient(db *pgx.Conn) Client {
	return &clientImpl{db: db}
}

func (c *clientImpl) RunInTransaction(ctx context.Context, fn func(context.Context, Transaction) error) error {
	// Start transaction
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to start sql transaction: %v", err)
	}

	// Run function
	err = fn(ctx, newTransaction(tx))
	if err != nil {
		// TODO check err
		tx.Rollback(ctx)
		return fmt.Errorf("sql transaction failed: %v", err)
	}

	// TODO support read-only transactions
	// TODO check err
	tx.Commit(ctx)

	return nil
}
