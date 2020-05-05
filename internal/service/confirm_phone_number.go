package service

import (
	"context"

	"github.com/AlpacaLabs/account-confirmation/internal/db"
	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
)

func (s Service) ConfirmPhoneNumber(ctx context.Context, request *confirmationV1.ConfirmPhoneNumberRequest) (*confirmationV1.ConfirmPhoneNumberResponse, error) {
	if err := s.dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
		// TODO lookup entity, mark as used, mark all codes for that phone number as confirmed
		return nil
	}); err != nil {
		return nil, err
	}

	return &confirmationV1.ConfirmPhoneNumberResponse{}, nil
}
