package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/AlpacaLabs/api-confirmation/internal/db/entities"

	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
)

const (
	TableForEmailCode = "email_address_confirmation_code"
	TableForPhoneCode = "phone_number_confirmation_code"
)

type Transaction interface {
	TransactionalOutbox

	CreateEmailCode(ctx context.Context, code confirmationV1.EmailAddressConfirmationCode) error
	CreatePhoneCode(ctx context.Context, code confirmationV1.PhoneNumberConfirmationCode) error

	GetEmailCodeByID(ctx context.Context, codeID string) (*confirmationV1.EmailAddressConfirmationCode, error)
	GetPhoneCodeByID(ctx context.Context, codeID string) (*confirmationV1.PhoneNumberConfirmationCode, error)
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

	queryTemplate := `
INSERT INTO %s(
  id, code, created_at, expires_at, stale, used, email_address_id
) 
VALUES($1, $2, $3, $4, $5, $6, $7)
`

	query := fmt.Sprintf(queryTemplate, TableForEmailCode)
	_, err := tx.tx.Exec(ctx, query, c.ID, c.Code, c.CreatedAt, c.ExpiresAt, c.Stale, c.Used, c.EmailAddressID)

	return err
}

func (tx *txImpl) CreatePhoneCode(ctx context.Context, in confirmationV1.PhoneNumberConfirmationCode) error {
	c := entities.NewPhoneNumberConfirmationCodeFromProtobuf(in)

	queryTemplate := `
INSERT INTO %s(
  id, code, created_at, expires_at, stale, used, phone_number_id
) 
VALUES($1, $2, $3, $4, $5, $6, $7)
`

	query := fmt.Sprintf(queryTemplate, TableForPhoneCode)
	_, err := tx.tx.Exec(ctx, query, c.ID, c.Code, c.CreatedAt, c.ExpiresAt, c.Stale, c.Used, c.PhoneNumberID)

	return err
}

func (tx *txImpl) GetEmailCodeByID(ctx context.Context, codeID string) (*confirmationV1.EmailAddressConfirmationCode, error) {
	queryTemplate := `
SELECT id, code, created_at, expires_at, stale, used, email_address_id
 FROM %s
 WHERE id = $1
`

	query := fmt.Sprintf(queryTemplate, TableForEmailCode)
	var c entities.EmailAddressConfirmationCode
	row := tx.tx.QueryRow(ctx, query, codeID)

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.EmailAddressID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) GetPhoneCodeByID(ctx context.Context, codeID string) (*confirmationV1.PhoneNumberConfirmationCode, error) {
	queryTemplate := `
SELECT id, code, created_at, expires_at, stale, used, phone_number_id
 FROM %s
 WHERE id = $1
`

	query := fmt.Sprintf(queryTemplate, TableForPhoneCode)
	var c entities.PhoneNumberConfirmationCode
	row := tx.tx.QueryRow(ctx, query, codeID)

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.PhoneNumberID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) GetEmailCode(ctx context.Context, code confirmationV1.ConfirmEmailAddressRequest) (*confirmationV1.EmailAddressConfirmationCode, error) {
	queryTemplate := `
SELECT id, code, created_at, expires_at, stale, used, email_address_id
 FROM %s
 WHERE code = $1
 AND email_address_id = $2
 AND stale = FALSE
 AND used = FALSE
 AND expires_at > $3
`

	query := fmt.Sprintf(queryTemplate, TableForEmailCode)
	var c entities.EmailAddressConfirmationCode
	row := tx.tx.QueryRow(ctx, query, code.Code, code.EmailAddressId, time.Now())

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.EmailAddressID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) GetPhoneCode(ctx context.Context, code confirmationV1.ConfirmPhoneNumberRequest) (*confirmationV1.PhoneNumberConfirmationCode, error) {
	queryTemplate := `
SELECT id, code, created_at, expires_at, stale, used, phone_number_id
 FROM %s
 WHERE code = $1
 AND phone_number_id = $2
 AND stale = FALSE
 AND used = FALSE
 AND expires_at > $3
`

	query := fmt.Sprintf(queryTemplate, TableForPhoneCode)
	var c entities.PhoneNumberConfirmationCode
	row := tx.tx.QueryRow(ctx, query, code.Code, code.PhoneNumberId, time.Now())

	err := row.Scan(&c.ID, &c.Code, &c.CreatedAt, &c.ExpiresAt, &c.Stale, &c.Used, &c.PhoneNumberID)
	if err != nil {
		return nil, err
	}

	return c.ToProtobuf(), nil
}

func (tx *txImpl) MarkEmailCodeAsUsed(ctx context.Context, codeID string) error {
	queryTemplate := `
UPDATE %s 
 SET used=TRUE, stale=TRUE 
 WHERE id=$1
`
	query := fmt.Sprintf(queryTemplate, TableForEmailCode)

	_, err := tx.tx.Exec(ctx, query, codeID)
	return err
}

func (tx *txImpl) MarkPhoneCodeAsUsed(ctx context.Context, codeID string) error {
	queryTemplate := `
UPDATE %s 
 SET used=TRUE, stale=TRUE 
 WHERE id=$1
`

	query := fmt.Sprintf(queryTemplate, TableForPhoneCode)
	_, err := tx.tx.Exec(ctx, query, codeID)
	return err
}

func (tx *txImpl) MarkEmailCodesAsStale(ctx context.Context, emailID string) error {
	queryTemplate := `
UPDATE %s 
 SET used=TRUE, stale=TRUE 
 WHERE email_address_id=$1
`
	query := fmt.Sprintf(queryTemplate, TableForEmailCode)

	_, err := tx.tx.Exec(ctx, query, emailID)
	return err
}

func (tx *txImpl) MarkPhoneCodesAsStale(ctx context.Context, phoneID string) error {
	queryTemplate := `
UPDATE %s 
 SET used=TRUE, stale=TRUE 
 WHERE phone_number_id=$1
`

	query := fmt.Sprintf(queryTemplate, TableForPhoneCode)
	_, err := tx.tx.Exec(ctx, query, phoneID)
	return err
}
