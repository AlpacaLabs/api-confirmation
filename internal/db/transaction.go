package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/AlpacaLabs/api-account-confirmation/internal/db/entities"

	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"

	"github.com/golang-sql/sqlexp"
)

type Transaction interface {
	CreateEmailConfirmationCode(ctx context.Context, code confirmationV1.EmailAddressConfirmationCode) error
	CreatePhoneConfirmationCode(ctx context.Context, code confirmationV1.PhoneNumberConfirmationCode) error

	CreateTxobForEmailCode(ctx context.Context, code confirmationV1.EmailAddressConfirmationCode) error
	CreateTxobForPhoneCode(ctx context.Context, code confirmationV1.PhoneNumberConfirmationCode) error

	GetEmailCode(ctx context.Context, code confirmationV1.ConfirmEmailAddressRequest) (*confirmationV1.EmailAddressConfirmationCode, error)
	GetPhoneCode(ctx context.Context, code confirmationV1.ConfirmPhoneNumberRequest) (*confirmationV1.PhoneNumberConfirmationCode, error)

	MarkEmailCodeAsUsed(ctx context.Context, codeID string) error
	MarkPhoneCodeAsUsed(ctx context.Context, codeID string) error

	MarkEmailCodesAsStale(ctx context.Context, emailID string) error
	MarkPhoneCodesAsStale(ctx context.Context, phoneID string) error
}

type txImpl struct {
	tx *sql.Tx
}

func (tx *txImpl) CreateEmailConfirmationCode(ctx context.Context, in confirmationV1.EmailAddressConfirmationCode) error {
	var q sqlexp.Querier
	q = tx.tx

	c := entities.NewEmailAddressConfirmationCodeFromProtobuf(in)

	query := `
INSERT INTO email_address_confirmation_code(
  id, code, created_timestamp, expiration_timestamp, stale, used, email_address_id
) 
VALUES($1, $2, $3, $4, $5, $6, $7)
`

	_, err := q.ExecContext(ctx, query, c.ID, c.Code, c.CreatedAt, c.ExpiresAt, c.Stale, c.Used, c.EmailAddressID)

	return err
}

func (tx *txImpl) CreatePhoneConfirmationCode(ctx context.Context, in confirmationV1.PhoneNumberConfirmationCode) error {
	var q sqlexp.Querier
	q = tx.tx

	c := entities.NewPhoneNumberConfirmationCodeFromProtobuf(in)

	query := `
INSERT INTO phone_number_confirmation_code(
  id, code, created_timestamp, expiration_timestamp, stale, used, phone_number_id
) 
VALUES($1, $2, $3, $4, $5, $6, $7)
`

	_, err := q.ExecContext(ctx, query, c.ID, c.Code, c.CreatedAt, c.ExpiresAt, c.Stale, c.Used, c.PhoneNumberID)

	return err
}

func (tx *txImpl) CreateTxobForEmailCode(ctx context.Context, code confirmationV1.EmailAddressConfirmationCode) error {
	var q sqlexp.Querier
	q = tx.tx

	query := `
INSERT INTO email_address_confirmation_code_txob(
  code_id, sent
) 
VALUES($1, FALSE)
`

	_, err := q.ExecContext(ctx, query, code.Id)

	return err
}

func (tx *txImpl) CreateTxobForPhoneCode(ctx context.Context, code confirmationV1.PhoneNumberConfirmationCode) error {
	var q sqlexp.Querier
	q = tx.tx

	query := `
INSERT INTO phone_number_confirmation_code_txob(
  code_id, sent
) 
VALUES($1, FALSE)
`

	_, err := q.ExecContext(ctx, query, code.Id)

	return err
}

func (tx *txImpl) GetEmailCode(ctx context.Context, code confirmationV1.ConfirmEmailAddressRequest) (*confirmationV1.EmailAddressConfirmationCode, error) {
	var q sqlexp.Querier
	q = tx.tx

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
	row := q.QueryRowContext(ctx, query, code.Code, code.EmailAddressId, time.Now())

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.EmailAddressID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) GetPhoneCode(ctx context.Context, code confirmationV1.ConfirmPhoneNumberRequest) (*confirmationV1.PhoneNumberConfirmationCode, error) {
	var q sqlexp.Querier
	q = tx.tx

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
	row := q.QueryRowContext(ctx, query, code.Code, code.PhoneNumberId, time.Now())

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.PhoneNumberID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) MarkEmailCodeAsUsed(ctx context.Context, codeID string) error {
	var q sqlexp.Querier
	q = tx.tx

	query := `
UPDATE email_address_confirmation_code 
 SET used=TRUE, stale=TRUE 
 WHERE id=$1
`

	_, err := q.ExecContext(ctx, query, codeID)
	return err
}

func (tx *txImpl) MarkPhoneCodeAsUsed(ctx context.Context, codeID string) error {
	var q sqlexp.Querier
	q = tx.tx

	query := `
UPDATE phone_number_confirmation_code 
 SET used=TRUE, stale=TRUE 
 WHERE id=$1
`

	_, err := q.ExecContext(ctx, query, codeID)
	return err
}

func (tx *txImpl) MarkEmailCodesAsStale(ctx context.Context, emailID string) error {
	var q sqlexp.Querier
	q = tx.tx

	query := `
UPDATE email_address_confirmation_code 
 SET used=TRUE, stale=TRUE 
 WHERE email_address_id=$1
`

	_, err := q.ExecContext(ctx, query, emailID)
	return err
}

func (tx *txImpl) MarkPhoneCodesAsStale(ctx context.Context, phoneID string) error {
	var q sqlexp.Querier
	q = tx.tx

	query := `
UPDATE phone_number_confirmation_code 
 SET used=TRUE, stale=TRUE 
 WHERE phone_number_id=$1
`

	_, err := q.ExecContext(ctx, query, phoneID)
	return err
}
