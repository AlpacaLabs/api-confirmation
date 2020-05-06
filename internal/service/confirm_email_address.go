package service

import (
	"context"

	"github.com/AlpacaLabs/api-account-confirmation/internal/db"
	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
)

func (s Service) ConfirmEmailAddress(ctx context.Context, request *confirmationV1.ConfirmEmailAddressRequest) (*confirmationV1.ConfirmEmailAddressResponse, error) {
	if err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		// TODO lookup entity, mark as used, mark all codes for that email as confirmed
		return nil
	}); err != nil {
		return nil, err
	}

	return &confirmationV1.ConfirmEmailAddressResponse{}, nil
}
