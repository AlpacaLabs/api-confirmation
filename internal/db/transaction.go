package db

import (
	"context"
	"database/sql"

	"github.com/AlpacaLabs/account-confirmation/internal/db/entities"

	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"

	"github.com/golang-sql/sqlexp"
)

type Transaction interface {
	CreateEmailConfirmationCode(ctx context.Context, code confirmationV1.EmailAddressConfirmationCode) error
	CreatePhoneConfirmationCode(ctx context.Context, code confirmationV1.PhoneNumberConfirmationCode) error
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
