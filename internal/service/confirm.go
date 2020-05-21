package service

import (
	"context"
	"fmt"

	"github.com/AlpacaLabs/go-kontext"

	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"

	"github.com/AlpacaLabs/api-confirmation/internal/db"
	"github.com/AlpacaLabs/api-confirmation/internal/db/entities"
	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
)

func (s Service) ConfirmEmailAddress(ctx context.Context, request *confirmationV1.ConfirmEmailAddressRequest) (*confirmationV1.ConfirmEmailAddressResponse, error) {
	funcName := "ConfirmEmailAddress"
	transactionalOutboxTable := db.TableForConfirmEmailAddressRequest

	if err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		// Verify the code is valid
		codeEntity, err := tx.GetEmailCode(ctx, *request)
		if err != nil {
			return err
		}

		// Mark code as used
		if err := tx.MarkEmailCodeAsUsed(ctx, codeEntity.Id); err != nil {
			return err
		}

		emailAddressID := codeEntity.EmailAddressId

		// Mark all codes for that email as stale
		if err := tx.MarkEmailCodesAsStale(ctx, emailAddressID); err != nil {
			return err
		}

		// Create the protocol buffer that the Message Relay process
		// will use as a payload in a Kafka topic.
		payload := &accountV1.ConfirmEmailAddressRequest{
			EmailAddressId: emailAddressID,
		}

		traceInfo := kontext.GetTraceInfo(ctx)

		// Create the event entity that will be persisted to the transactional outbox
		event, err := entities.NewEvent(traceInfo, request, payload)
		if err != nil {
			return fmt.Errorf("failed to create event in %s: %w", funcName, err)
		}

		// Persist the event to the transactional outbox
		return tx.CreateEvent(ctx, event, transactionalOutboxTable)
	}); err != nil {
		return nil, err
	}

	return &confirmationV1.ConfirmEmailAddressResponse{}, nil
}

func (s Service) ConfirmPhoneNumber(ctx context.Context, request *confirmationV1.ConfirmPhoneNumberRequest) (*confirmationV1.ConfirmPhoneNumberResponse, error) {
	funcName := "ConfirmPhoneNumber"
	transactionalOutboxTable := db.TableForConfirmPhoneNumberRequest

	if err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		// Verify the code is valid
		codeEntity, err := tx.GetPhoneCode(ctx, *request)
		if err != nil {
			return err
		}

		// Mark code as used
		if err := tx.MarkPhoneCodeAsUsed(ctx, codeEntity.Id); err != nil {
			return err
		}

		phoneNumberID := codeEntity.PhoneNumberId

		// Mark all codes for that phone as stale
		if err := tx.MarkPhoneCodesAsStale(ctx, phoneNumberID); err != nil {
			return err
		}

		// Create the protocol buffer that the Message Relay process
		// will use as a payload in a Kafka topic.
		payload := &accountV1.ConfirmPhoneNumberRequest{
			PhoneNumberId: phoneNumberID,
		}

		traceInfo := kontext.GetTraceInfo(ctx)

		// Create the event entity that will be persisted to the transactional outbox
		event, err := entities.NewEvent(traceInfo, request, payload)
		if err != nil {
			return fmt.Errorf("failed to create event in %s: %w", funcName, err)
		}

		// Persist the event to the transactional outbox
		return tx.CreateEvent(ctx, event, transactionalOutboxTable)
	}); err != nil {
		return nil, err
	}

	return &confirmationV1.ConfirmPhoneNumberResponse{}, nil
}
