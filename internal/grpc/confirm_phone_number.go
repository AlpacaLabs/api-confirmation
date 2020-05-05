package grpc

import (
	"context"

	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
)

func (s Server) ConfirmPhoneNumber(ctx context.Context, request *confirmationV1.ConfirmPhoneNumberRequest) (*confirmationV1.ConfirmPhoneNumberResponse, error) {
	return s.service.ConfirmPhoneNumber(ctx, request)
}
