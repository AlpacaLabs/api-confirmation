package service

import (
	"context"

	"github.com/AlpacaLabs/api-confirmation/internal/db"
	"github.com/AlpacaLabs/api-confirmation/internal/db/entities"
	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
)

func (s Service) ConfirmEmailAddress(ctx context.Context, request *confirmationV1.ConfirmEmailAddressRequest) (*confirmationV1.ConfirmEmailAddressResponse, error) {
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

		// Write to transactional outbox so a message relay
		// can mark the email address as confirmed in the
		// Account Service.
		return tx.CreateTxobForEmailConfirmation(ctx, entities.NewConfirmEmailEvent(ctx, emailAddressID))
	}); err != nil {
		return nil, err
	}

	return &confirmationV1.ConfirmEmailAddressResponse{}, nil
}

func (s Service) ConfirmPhoneNumber(ctx context.Context, request *confirmationV1.ConfirmPhoneNumberRequest) (*confirmationV1.ConfirmPhoneNumberResponse, error) {
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

		// Write to transactional outbox so a message relay
		// can mark the phone number as confirmed in the
		// Account Service.
		return tx.CreateTxobForPhoneConfirmation(ctx, entities.NewConfirmPhoneEvent(ctx, phoneNumberID))
	}); err != nil {
		return nil, err
	}

	return &confirmationV1.ConfirmPhoneNumberResponse{}, nil
}
