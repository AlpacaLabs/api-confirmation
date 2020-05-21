package db

import (
	"context"
	"fmt"

	"github.com/AlpacaLabs/api-confirmation/internal/db/entities"
	"github.com/jackc/pgx/v4"
)

const (
	// TableForSendEmailRequest is the name of the transactional outbox (database table)
	// from which we read "jobs" or "events" that need to get sent to a message broker.
	TableForSendEmailRequest = "txob_send_email_request"
	TableForSendSmsRequest   = "txob_send_sms_request"

	TableForConfirmEmailAddressRequest = "txob_confirm_email_address_request"
	TableForConfirmPhoneNumberRequest  = "txob_confirm_phone_number_request"
)

type TransactionalOutbox interface {
	ReadEvent(ctx context.Context, table string) (e *entities.Event, err error)
	CreateEvent(ctx context.Context, e entities.Event, table string) error
	MarkEventAsSent(ctx context.Context, eventID string, table string) error
}

type outboxImpl struct {
	tx pgx.Tx
}

func (tx *outboxImpl) ReadEvent(ctx context.Context, table string) (*entities.Event, error) {
	queryTemplate := `
SELECT
  event_id, trace_id, sampled, sent, catalyst, payload
  FROM %s
  WHERE sent = FALSE
  LIMIT 1
`

	query := fmt.Sprintf(queryTemplate, table)

	row := tx.tx.QueryRow(ctx, query)

	e := &entities.Event{}

	if err := row.Scan(&e.EventId, &e.TraceId, &e.Sampled, &e.Sent, &e.Catalyst, &e.Payload); err != nil {
		return e, err
	}

	return e, nil
}

func (tx *outboxImpl) CreateEvent(ctx context.Context, e entities.Event, table string) error {
	queryTemplate := `
INSERT INTO %s(
  event_id, trace_id, sampled, sent, catalyst, payload
) 
VALUES($1, $2, $3, $4, $5, $6)
`

	query := fmt.Sprintf(queryTemplate, table)

	if _, err := tx.tx.Exec(ctx, query, e.EventId, e.TraceId, e.Sampled, e.Sent, e.Catalyst, e.Payload); err != nil {
		return err
	}

	return nil
}

func (tx *outboxImpl) MarkEventAsSent(ctx context.Context, eventID string, table string) error {
	queryTemplate := `
UPDATE %s
  SET sent = TRUE
  WHERE event_id = $1
`
	query := fmt.Sprintf(queryTemplate, table)
	_, err := tx.tx.Exec(ctx, query, eventID)
	return err
}
