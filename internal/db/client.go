package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/sirupsen/logrus"
)

type TxOption string

const ReadOnly = TxOption("read-only")

type Client interface {
	RunInTransaction(ctx context.Context, fn func(context.Context, Transaction) error, options ...TxOption) error
}

type clientImpl struct {
	db *pgx.Conn
}

func NewClient(db *pgx.Conn) Client {
	return &clientImpl{db: db}
}

func (c *clientImpl) RunInTransaction(ctx context.Context, fn func(context.Context, Transaction) error, options ...TxOption) error {
	var readOnly bool
	for _, o := range options {
		if o == ReadOnly {
			readOnly = true
		}
	}

	opts := pgx.TxOptions{}
	if readOnly {
		opts.AccessMode = pgx.ReadOnly
	}

	// Start transaction
	tx, err := c.db.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed to start sql transaction: %v", err)
	}

	defer func() {
		if !readOnly {
			if err := tx.Rollback(ctx); err != nil {
				logrus.Errorf("failed to rollback transaction: %v", err)
			}
		}
	}()

	// Run function
	err = fn(ctx, newTransaction(tx))
	if err != nil {
		return fmt.Errorf("sql transaction failed: %v", err)
	}

	if !readOnly {
		return tx.Commit(ctx)
	}

	return nil
}
