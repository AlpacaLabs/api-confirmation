package service

import (
	"context"

	"github.com/AlpacaLabs/api-account-confirmation/internal/db"
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

		// Mark all codes for that email as stale
		if err := tx.MarkEmailCodesAsStale(ctx, codeEntity.EmailAddressId); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &confirmationV1.ConfirmEmailAddressResponse{}, nil
}
