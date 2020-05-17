package db

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/AlpacaLabs/api-confirmation/internal/db/entities"

	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
)

type Transaction interface {
	TransactionalOutbox

	CreateEmailCode(ctx context.Context, code confirmationV1.EmailAddressConfirmationCode) error
	CreatePhoneCode(ctx context.Context, code confirmationV1.PhoneNumberConfirmationCode) error

	GetEmailCode(ctx context.Context, code confirmationV1.ConfirmEmailAddressRequest) (*confirmationV1.EmailAddressConfirmationCode, error)
	GetPhoneCode(ctx context.Context, code confirmationV1.ConfirmPhoneNumberRequest) (*confirmationV1.PhoneNumberConfirmationCode, error)

	MarkEmailCodeAsUsed(ctx context.Context, codeID string) error
	MarkPhoneCodeAsUsed(ctx context.Context, codeID string) error

	MarkEmailCodesAsStale(ctx context.Context, emailID string) error
	MarkPhoneCodesAsStale(ctx context.Context, phoneID string) error
}

type txImpl struct {
	tx pgx.Tx
	outboxImpl
}

func newTransaction(tx pgx.Tx) Transaction {
	return &txImpl{
		tx: tx,
		outboxImpl: outboxImpl{
			tx: tx,
		},
	}
}

func (tx *txImpl) CreateEmailCode(ctx context.Context, in confirmationV1.EmailAddressConfirmationCode) error {
	c := entities.NewEmailAddressConfirmationCodeFromProtobuf(in)

	query := `
INSERT INTO email_address_confirmation_code(
  id, code, created_timestamp, expiration_timestamp, stale, used, email_address_id
) 
VALUES($1, $2, $3, $4, $5, $6, $7)
`

	_, err := tx.tx.Exec(ctx, query, c.ID, c.Code, c.CreatedAt, c.ExpiresAt, c.Stale, c.Used, c.EmailAddressID)

	return err
}

func (tx *txImpl) CreatePhoneCode(ctx context.Context, in confirmationV1.PhoneNumberConfirmationCode) error {
	c := entities.NewPhoneNumberConfirmationCodeFromProtobuf(in)

	query := `
INSERT INTO phone_number_confirmation_code(
  id, code, created_timestamp, expiration_timestamp, stale, used, phone_number_id
) 
VALUES($1, $2, $3, $4, $5, $6, $7)
`

	_, err := tx.tx.Exec(ctx, query, c.ID, c.Code, c.CreatedAt, c.ExpiresAt, c.Stale, c.Used, c.PhoneNumberID)

	return err
}

func (tx *txImpl) GetEmailCode(ctx context.Context, code confirmationV1.ConfirmEmailAddressRequest) (*confirmationV1.EmailAddressConfirmationCode, error) {
	query := `
SELECT id, code, creation_timestamp, expiration_timestamp, stale, used, email_address_id
 FROM email_address_confirmation_code
 WHERE code = $1
 AND email_address_id = $2
 AND stale = FALSE
 AND used = FALSE
 AND expiration_timestamp > $3
`

	var c entities.EmailAddressConfirmationCode
	row := tx.tx.QueryRow(ctx, query, code.Code, code.EmailAddressId, time.Now())

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.EmailAddressID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) GetPhoneCode(ctx context.Context, code confirmationV1.ConfirmPhoneNumberRequest) (*confirmationV1.PhoneNumberConfirmationCode, error) {
	query := `
SELECT id, code, creation_timestamp, expiration_timestamp, stale, used, phone_number_id
 FROM phone_number_confirmation_code
 WHERE code = $1
 AND phone_number_id = $2
 AND stale = FALSE
 AND used = FALSE
 AND expiration_timestamp > $3
`

	var c entities.PhoneNumberConfirmationCode
	row := tx.tx.QueryRow(ctx, query, code.Code, code.PhoneNumberId, time.Now())

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.PhoneNumberID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) MarkEmailCodeAsUsed(ctx context.Context, codeID string) error {
	query := `
UPDATE email_address_confirmation_code 
 SET used=TRUE, stale=TRUE 
 WHERE id=$1
`

	_, err := tx.tx.Exec(ctx, query, codeID)
	return err
}

func (tx *txImpl) MarkPhoneCodeAsUsed(ctx context.Context, codeID string) error {
	query := `
UPDATE phone_number_confirmation_code 
 SET used=TRUE, stale=TRUE 
 WHERE id=$1
`

	_, err := tx.tx.Exec(ctx, query, codeID)
	return err
}

func (tx *txImpl) MarkEmailCodesAsStale(ctx context.Context, emailID string) error {
	query := `
UPDATE email_address_confirmation_code 
 SET used=TRUE, stale=TRUE 
 WHERE email_address_id=$1
`

	_, err := tx.tx.Exec(ctx, query, emailID)
	return err
}

func (tx *txImpl) MarkPhoneCodesAsStale(ctx context.Context, phoneID string) error {
	query := `
UPDATE phone_number_confirmation_code 
 SET used=TRUE, stale=TRUE 
 WHERE phone_number_id=$1
`

	_, err := tx.tx.Exec(ctx, query, phoneID)
	return err
}
