package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Client interface {
	RunInTransaction(ctx context.Context, fn func(context.Context, Transaction) error) error
}

type clientImpl struct {
	db *sql.DB
}

func NewClient(db *sql.DB) Client {
	return &clientImpl{db: db}
}

func (c *clientImpl) RunInTransaction(ctx context.Context, fn func(context.Context, Transaction) error) error {
	// Start transaction
	tx, err := c.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start sql transaction: %v", err)
	}

	// Run function
	err = fn(ctx, &txImpl{
		tx: tx,
	})
	if err != nil {
		// TODO check err
		tx.Rollback()
		return fmt.Errorf("sql transaction failed: %v", err)
	}

	// TODO support read-only transactions
	// TODO check err
	tx.Commit()

	return nil
}
