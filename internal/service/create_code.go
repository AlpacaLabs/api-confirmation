package service

import (
	"context"
	"fmt"
	"time"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	hermesV1 "github.com/AlpacaLabs/protorepo-hermes-go/alpacalabs/hermes/v1"

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
	funcName := "CreateEmailAddressConfirmationCode"
	transactionalOutboxTable := db.TableForSendEmailRequest
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

		// Create the protocol buffer that the Message Relay process
		// will use as a payload in a Kafka topic.
		payload, err := s.buildEmail(ctx, c.EmailAddressId)
		if err != nil {
			return err
		}

		// Create the event entity that will be persisted to the transactional outbox
		event, err := entities.NewEvent(ctx, request, payload)
		if err != nil {
			return fmt.Errorf("failed to create event in %s: %w", funcName, err)
		}

		// Persist the event to the transactional outbox
		return tx.CreateEvent(ctx, event, transactionalOutboxTable)
	})
}

func (s Service) CreatePhoneNumberConfirmationCode(ctx context.Context, traceInfo eventV1.TraceInfo, request *confirmationV1.CreatePhoneNumberConfirmationCodeRequest) error {
	funcName := "CreatePhoneNumberConfirmationCode"
	transactionalOutboxTable := db.TableForSendSmsRequest
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

		// Create the protocol buffer that the Message Relay process
		// will use as a payload in a Kafka topic.
		payload, err := s.buildSms(ctx, c.PhoneNumberId)
		if err != nil {
			return err
		}

		// Create the event entity that will be persisted to the transactional outbox
		event, err := entities.NewEvent(ctx, request, payload)
		if err != nil {
			return fmt.Errorf("failed to create event in %s: %w", funcName, err)
		}

		// Persist the event to the transactional outbox
		return tx.CreateEvent(ctx, event, transactionalOutboxTable)
	})
}

func (s Service) buildEmail(ctx context.Context, emailAddressID string) (*hermesV1.SendEmailRequest, error) {
	var emailAddress string
	accountClient := accountV1.NewAccountServiceClient(s.accountConn)
	if emailRes, err := accountClient.GetEmailAddress(ctx, &accountV1.GetEmailAddressRequest{
		Id: emailAddressID,
	}); err != nil {
		return nil, fmt.Errorf("failed to verify email address existence: %w", err)
	} else {
		emailAddress = emailRes.EmailAddress.EmailAddress
	}

	pb := &hermesV1.SendEmailRequest{
		// TODO build actual email body to let them know they should confirm their account
		Email: &hermesV1.Email{
			Subject: "Email Confirmation",
			Body: &hermesV1.Body{
				Intros: []string{
					"Please click the link to confirm your email",
				},
				Actions: []*hermesV1.Action{
					{
						Instructions: "Click to Confirm",
						Button: &hermesV1.Button{
							Color:     "",
							TextColor: "",
							Text:      "Click to Confirm",
							Link:      "",
						},
					},
				},
				Signature: "Welcome",
			},
			To: []*hermesV1.Recipient{
				{
					Email: emailAddress,
				},
			},
		},
	}
	return pb, nil
}

func (s Service) buildSms(ctx context.Context, phoneNumberID string) (*hermesV1.SendSmsRequest, error) {
	// TODO build request
	return &hermesV1.SendSmsRequest{
		To:      "",
		Message: "",
	}, nil
}
