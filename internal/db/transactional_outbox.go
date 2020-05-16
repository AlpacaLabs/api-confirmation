package db

import (
	"context"

	"github.com/AlpacaLabs/api-account-confirmation/internal/db/entities"
	"github.com/jackc/pgx/v4"
)

type TransactionalOutbox interface {
	CreateTxobForEmailCode(ctx context.Context, e entities.SendEmailEvent) error
	CreateTxobForPhoneCode(ctx context.Context, e entities.SendPhoneEvent) error
	ReadFromTxobForEmailCode(ctx context.Context, num int) ([]entities.SendEmailEvent, error)
	ReadFromTxobForPhoneCode(ctx context.Context, num int) ([]entities.SendPhoneEvent, error)
}

type outboxImpl struct {
	tx pgx.Tx
}

func (tx *outboxImpl) CreateTxobForEmailCode(ctx context.Context, e entities.SendEmailEvent) error {
	query := `
INSERT INTO email_address_confirmation_code_txob(
  event_id, code_id, sent
) 
VALUES($1, FALSE)
`

	_, err := tx.tx.Exec(ctx, query, e.EventID, e.CodeID, e.Sent)

	return err
}

func (tx *outboxImpl) CreateTxobForPhoneCode(ctx context.Context, e entities.SendPhoneEvent) error {
	query := `
INSERT INTO phone_number_confirmation_code_txob(
  event_id, code_id, sent
) 
VALUES($1, FALSE)
`

	_, err := tx.tx.Exec(ctx, query, e.EventID, e.CodeID, e.Sent)

	return err
}

func (tx *outboxImpl) ReadFromTxobForEmailCode(ctx context.Context, num int) ([]entities.SendEmailEvent, error) {
	query := `
SELECT event_id, code_id
  FROM email_address_confirmation_code_txob
  WHERE sent = FALSE
  LIMIT $1
`
	rows, err := tx.tx.Query(ctx, query, num)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []entities.SendEmailEvent{}

	for rows.Next() {
		var e entities.SendEmailEvent
		if err := rows.Scan(&e.EventID, &e.CodeID); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (tx *outboxImpl) ReadFromTxobForPhoneCode(ctx context.Context, num int) ([]entities.SendPhoneEvent, error) {
	query := `
SELECT event_id, code_id
  FROM email_address_confirmation_code_txob
  WHERE sent = FALSE
  LIMIT $1
`
	rows, err := tx.tx.Query(ctx, query, num)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []entities.SendPhoneEvent{}

	for rows.Next() {
		var e entities.SendPhoneEvent
		if err := rows.Scan(&e.EventID, &e.CodeID); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}
