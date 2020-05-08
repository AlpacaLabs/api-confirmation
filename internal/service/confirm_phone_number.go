package service

import (
	"context"

	"github.com/AlpacaLabs/api-account-confirmation/internal/db"
	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
)

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

		// Mark all codes for that phone as stale
		if err := tx.MarkPhoneCodesAsStale(ctx, codeEntity.PhoneNumberId); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &confirmationV1.ConfirmPhoneNumberResponse{}, nil
}
