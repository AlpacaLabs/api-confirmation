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
	ReadEventForSendEmailRequest(ctx context.Context) (e entities.SendEmailEvent, err error)
	ReadEventForSendSmsRequest(ctx context.Context) (e entities.SendPhoneEvent, err error)
	CreateEventForSendEmailRequest(ctx context.Context, e entities.SendEmailEvent) error
	CreateEventForSendSmsRequest(ctx context.Context, e entities.SendPhoneEvent) error
	MarkAsSentSendEmailRequest(ctx context.Context, eventID string) error
	MarkAsSentSendSmsRequest(ctx context.Context, eventID string) error

	ReadEventForEmailConfirmation(ctx context.Context) (e entities.ConfirmEmailEvent, err error)
	ReadEventForPhoneConfirmation(ctx context.Context) (e entities.ConfirmPhoneEvent, err error)
	CreateEventForEmailConfirmation(ctx context.Context, e entities.ConfirmEmailEvent) error
	CreateEventForPhoneConfirmation(ctx context.Context, e entities.ConfirmPhoneEvent) error
	MarkAsSentConfirmEmailEvent(ctx context.Context, eventID string) error
	MarkAsSentConfirmPhoneEvent(ctx context.Context, eventID string) error
}

type outboxImpl struct {
	tx pgx.Tx
}

func (tx *outboxImpl) ReadEventForSendEmailRequest(ctx context.Context) (e entities.SendEmailEvent, err error) {
	queryTemplate := `
SELECT event_id, trace_id, sampled, code_id
  FROM %s
  WHERE sent = FALSE
  LIMIT 1
`

	query := fmt.Sprintf(queryTemplate, TableForSendEmailRequest)

	row := tx.tx.QueryRow(ctx, query)

	if err := row.Scan(&e.EventId, &e.TraceId, &e.Sampled, &e.CodeID); err != nil {
		return e, err
	}

	return e, nil
}

func (tx *outboxImpl) ReadEventForSendSmsRequest(ctx context.Context) (e entities.SendPhoneEvent, err error) {
	queryTemplate := `
SELECT event_id, trace_id, sampled, code_id
  FROM %s
  WHERE sent = FALSE
  LIMIT 1
`

	query := fmt.Sprintf(queryTemplate, TableForSendSmsRequest)

	row := tx.tx.QueryRow(ctx, query)

	if err := row.Scan(&e.EventId, &e.TraceId, &e.Sampled, &e.CodeID); err != nil {
		return e, err
	}

	return e, nil
}

func (tx *outboxImpl) CreateEventForSendEmailRequest(ctx context.Context, e entities.SendEmailEvent) error {
	queryTemplate := `
INSERT INTO %s(
  event_id, trace_id, sampled, sent, code_id
) 
VALUES($1, $2, $3, FALSE, $4)
`

	query := fmt.Sprintf(queryTemplate, TableForSendEmailRequest)

	_, err := tx.tx.Exec(ctx, query, e.EventId, e.TraceId, e.Sampled, e.CodeID)

	return err
}

func (tx *outboxImpl) CreateEventForSendSmsRequest(ctx context.Context, e entities.SendPhoneEvent) error {
	queryTemplate := `
INSERT INTO %s(
  event_id, trace_id, sampled, sent, code_id
) 
VALUES($1, $2, $3, FALSE, $4)
`

	query := fmt.Sprintf(queryTemplate, TableForSendSmsRequest)

	_, err := tx.tx.Exec(ctx, query, e.EventId, e.TraceId, e.Sampled, e.CodeID)

	return err
}

func (tx *outboxImpl) MarkAsSentSendEmailRequest(ctx context.Context, eventID string) error {
	queryTemplate := `
UPDATE %s
  SET sent = TRUE
  WHERE event_id = $1
`
	query := fmt.Sprintf(queryTemplate, TableForSendEmailRequest)
	_, err := tx.tx.Exec(ctx, query, eventID)
	return err
}

func (tx *outboxImpl) MarkAsSentSendSmsRequest(ctx context.Context, eventID string) error {
	queryTemplate := `
UPDATE %s
  SET sent = TRUE
  WHERE event_id = $1
`
	query := fmt.Sprintf(queryTemplate, TableForSendSmsRequest)
	_, err := tx.tx.Exec(ctx, query, eventID)
	return err
}

func (tx *outboxImpl) ReadEventForEmailConfirmation(ctx context.Context) (e entities.ConfirmEmailEvent, err error) {
	queryTemplate := `
SELECT event_id, trace_id, sampled, email_address_id
  FROM %s
  WHERE sent = FALSE
  LIMIT 1
`
	query := fmt.Sprintf(queryTemplate, TableForConfirmEmailAddressRequest)

	row := tx.tx.QueryRow(ctx, query)

	if err := row.Scan(&e.EventId, &e.TraceId, &e.Sampled, &e.EmailAddressID); err != nil {
		return e, err
	}

	return e, nil
}

func (tx *outboxImpl) ReadEventForPhoneConfirmation(ctx context.Context) (e entities.ConfirmPhoneEvent, err error) {
	queryTemplate := `
SELECT event_id, trace_id, sampled, phone_number_id
  FROM %s
  WHERE sent = FALSE
  LIMIT 1
`

	query := fmt.Sprintf(queryTemplate, TableForConfirmPhoneNumberRequest)

	row := tx.tx.QueryRow(ctx, query)

	if err := row.Scan(&e.EventId, &e.TraceId, &e.Sampled, &e.PhoneNumberID); err != nil {
		return e, err
	}

	return e, err
}

func (tx *outboxImpl) CreateEventForEmailConfirmation(ctx context.Context, e entities.ConfirmEmailEvent) error {
	queryTemplate := `
INSERT INTO %s(
  event_id, trace_id, sampled, sent, email_address_id
) 
VALUES($1, $2, $3, FALSE, $4)
`

	query := fmt.Sprintf(queryTemplate, TableForConfirmEmailAddressRequest)

	_, err := tx.tx.Exec(ctx, query, e.EventId, e.TraceId, e.Sampled, e.EmailAddressID)

	return err
}

func (tx *outboxImpl) CreateEventForPhoneConfirmation(ctx context.Context, e entities.ConfirmPhoneEvent) error {
	queryTemplate := `
INSERT INTO %s(
  event_id, trace_id, sampled, sent, phone_number_id
) 
VALUES($1, $2, $3, FALSE, $4)
`

	query := fmt.Sprintf(queryTemplate, TableForConfirmPhoneNumberRequest)

	_, err := tx.tx.Exec(ctx, query, e.EventId, e.TraceId, e.Sampled, e.PhoneNumberID)

	return err
}

func (tx *outboxImpl) MarkAsSentConfirmEmailEvent(ctx context.Context, eventID string) error {
	queryTemplate := `
UPDATE %s
  SET sent = TRUE
  WHERE event_id = $1
`
	query := fmt.Sprintf(queryTemplate, TableForConfirmEmailAddressRequest)
	_, err := tx.tx.Exec(ctx, query, eventID)
	return err
}

func (tx *outboxImpl) MarkAsSentConfirmPhoneEvent(ctx context.Context, eventID string) error {
	queryTemplate := `
UPDATE %s
  SET sent = TRUE
  WHERE event_id = $1
`
	query := fmt.Sprintf(queryTemplate, TableForConfirmPhoneNumberRequest)
	_, err := tx.tx.Exec(ctx, query, eventID)
	return err
}
