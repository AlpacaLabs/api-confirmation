package service

import (
	"context"
	"time"

	"github.com/AlpacaLabs/api-account-confirmation/internal/db"
	clock "github.com/AlpacaLabs/go-timestamp"
	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
	"github.com/rs/xid"
)

const (
	codeLongevity = time.Hour * 24
)

func (s Service) CreateEmailAddressConfirmationCode(ctx context.Context, request *confirmationV1.CreateEmailAddressConfirmationCodeRequest) error {
	return s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		now := time.Now()
		c := confirmationV1.EmailAddressConfirmationCode{
			Id:             xid.New().String(),
			EmailAddressId: request.EmailAddressId,
			Code:           code.New(),
			CreatedAt:      clock.TimeToTimestamp(now),
			ExpiresAt:      clock.TimeToTimestamp(now.Add(codeLongevity)),
		}
		return tx.CreateEmailCode(ctx, c)

		// TODO also write TXOB record to hit Hermes
	})
}

func (s Service) CreatePhoneNumberConfirmationCode(ctx context.Context, request *confirmationV1.CreatePhoneNumberConfirmationCodeRequest) error {
	return s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		now := time.Now()
		c := confirmationV1.PhoneNumberConfirmationCode{
			Id:            xid.New().String(),
			PhoneNumberId: request.PhoneNumberId,
			Code:          code.New(),
			CreatedAt:     clock.TimeToTimestamp(now),
			ExpiresAt:     clock.TimeToTimestamp(now.Add(codeLongevity)),
		}
		return tx.CreatePhoneCode(ctx, c)

		// TODO also write TXOB record to hit Hermes
	})
}
