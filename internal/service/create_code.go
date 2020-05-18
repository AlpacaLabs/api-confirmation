package service

import (
	"context"
	"time"

	eventV1 "github.com/AlpacaLabs/protorepo-event-go/alpacalabs/event/v1"

	"github.com/AlpacaLabs/api-confirmation/internal/db/entities"
	code "github.com/AlpacaLabs/go-random-code"

	"github.com/AlpacaLabs/api-confirmation/internal/db"
	clock "github.com/AlpacaLabs/go-timestamp"
	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
	"github.com/rs/xid"
)

const (
	codeLongevity = time.Hour * 24
)

func (s Service) CreateEmailAddressConfirmationCode(ctx context.Context, traceInfo eventV1.TraceInfo, request *confirmationV1.CreateEmailAddressConfirmationCodeRequest) error {
	return s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		now := time.Now()
		c := confirmationV1.EmailAddressConfirmationCode{
			Id:             xid.New().String(),
			EmailAddressId: request.EmailAddressId,
			Code:           code.New(),
			CreatedAt:      clock.TimeToTimestamp(now),
			ExpiresAt:      clock.TimeToTimestamp(now.Add(codeLongevity)),
		}
		if err := tx.CreateEmailCode(ctx, c); err != nil {
			return err
		}

		return tx.CreateEventForSendEmailRequest(ctx, entities.NewSendEmailEvent(traceInfo, c.Id))
	})
}

func (s Service) CreatePhoneNumberConfirmationCode(ctx context.Context, traceInfo eventV1.TraceInfo, request *confirmationV1.CreatePhoneNumberConfirmationCodeRequest) error {
	return s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		now := time.Now()
		c := confirmationV1.PhoneNumberConfirmationCode{
			Id:            xid.New().String(),
			PhoneNumberId: request.PhoneNumberId,
			Code:          code.New(),
			CreatedAt:     clock.TimeToTimestamp(now),
			ExpiresAt:     clock.TimeToTimestamp(now.Add(codeLongevity)),
		}
		if err := tx.CreatePhoneCode(ctx, c); err != nil {
			return err
		}

		return tx.CreateEventForSendSmsRequest(ctx, entities.NewSendPhoneEvent(traceInfo, c.Id))
	})
}
