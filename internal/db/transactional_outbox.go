package db

import (
	"context"

	"github.com/AlpacaLabs/api-confirmation/internal/db/entities"
	"github.com/jackc/pgx/v4"
)

type TransactionalOutbox interface {
	CreateTxobForEmailCode(ctx context.Context, e entities.SendEmailEvent) error
	CreateTxobForPhoneCode(ctx context.Context, e entities.SendPhoneEvent) error
	ReadFromTxobForEmailCode(ctx context.Context, num int) ([]entities.SendEmailEvent, error)
	ReadFromTxobForPhoneCode(ctx context.Context, num int) ([]entities.SendPhoneEvent, error)

	CreateTxobForEmailConfirmation(ctx context.Context, e entities.ConfirmEmailEvent) error
	ReadFromTxobForEmailConfirmation(ctx context.Context, num int) ([]entities.ConfirmEmailEvent, error)

	CreateTxobForPhoneConfirmation(ctx context.Context, e entities.ConfirmPhoneEvent) error
	ReadFromTxobForPhoneConfirmation(ctx context.Context, num int) ([]entities.ConfirmPhoneEvent, error)
}

type outboxImpl struct {
	tx pgx.Tx
}

func (tx *outboxImpl) CreateTxobForEmailCode(ctx context.Context, e entities.SendEmailEvent) error {
	query := `
INSERT INTO email_address_confirmation_code_txob(
  event_id, trace_id, sampled, sent, code_id
) 
VALUES($1, $2, $3, FALSE, $4)
`

	_, err := tx.tx.Exec(ctx, query, e.EventId, e.TraceId, e.Sampled, e.CodeID)

	return err
}

func (tx *outboxImpl) CreateTxobForPhoneCode(ctx context.Context, e entities.SendPhoneEvent) error {
	query := `
INSERT INTO phone_number_confirmation_code_txob(
  event_id, trace_id, sampled, sent, code_id
) 
VALUES($1, $2, $3, FALSE, $4)
`

	_, err := tx.tx.Exec(ctx, query, e.EventId, e.TraceId, e.Sampled, e.CodeID)

	return err
}

func (tx *outboxImpl) ReadFromTxobForEmailCode(ctx context.Context, num int) ([]entities.SendEmailEvent, error) {
	query := `
SELECT event_id, trace_id, sampled, code_id
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
		if err := rows.Scan(&e.EventId, &e.TraceId, &e.Sampled, &e.CodeID); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (tx *outboxImpl) ReadFromTxobForPhoneCode(ctx context.Context, num int) ([]entities.SendPhoneEvent, error) {
	query := `
SELECT event_id, trace_id, sampled, code_id
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
		if err := rows.Scan(&e.EventId, &e.TraceId, &e.Sampled, &e.CodeID); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (tx *outboxImpl) ReadFromTxobForPhoneConfirmation(ctx context.Context, num int) ([]entities.ConfirmPhoneEvent, error) {
	query := `
SELECT event_id, trace_id, sampled, phone_number_id
  FROM phone_number_confirmation_txob
  WHERE sent = FALSE
  LIMIT $1
`
	rows, err := tx.tx.Query(ctx, query, num)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []entities.ConfirmPhoneEvent{}

	for rows.Next() {
		var e entities.ConfirmPhoneEvent
		if err := rows.Scan(&e.EventId, &e.TraceId, &e.Sampled, &e.PhoneNumberID); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (tx *outboxImpl) ReadFromTxobForEmailConfirmation(ctx context.Context, num int) ([]entities.ConfirmEmailEvent, error) {
	query := `
SELECT event_id, trace_id, sampled, email_address_id
  FROM email_address_confirmation_txob
  WHERE sent = FALSE
  LIMIT $1
`
	rows, err := tx.tx.Query(ctx, query, num)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	events := []entities.ConfirmEmailEvent{}

	for rows.Next() {
		var e entities.ConfirmEmailEvent
		if err := rows.Scan(&e.EventId, &e.TraceId, &e.Sampled, &e.EmailAddressID); err != nil {
			return nil, err
		}
		events = append(events, e)
	}

	return events, nil
}

func (tx *outboxImpl) CreateTxobForEmailConfirmation(ctx context.Context, e entities.ConfirmEmailEvent) error {
	query := `
INSERT INTO email_address_confirmation_txob(
  event_id, trace_id, sampled, sent, email_address_id
) 
VALUES($1, $2, $3, FALSE, $4)
`

	_, err := tx.tx.Exec(ctx, query, e.EventId, e.TraceId, e.Sampled, e.EmailAddressID)

	return err
}

func (tx *outboxImpl) CreateTxobForPhoneConfirmation(ctx context.Context, e entities.ConfirmPhoneEvent) error {
	query := `
INSERT INTO phone_number_confirmation_txob(
  event_id, trace_id, sampled, sent, phone_number_id
) 
VALUES($1, $2, $3, FALSE, $4)
`

	_, err := tx.tx.Exec(ctx, query, e.EventId, e.TraceId, e.Sampled, e.PhoneNumberID)

	return err
}
