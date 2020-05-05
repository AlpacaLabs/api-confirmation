package grpc

import (
	"context"

	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
)

func (s Server) ConfirmEmailAddress(ctx context.Context, request *confirmationV1.ConfirmEmailAddressRequest) (*confirmationV1.ConfirmEmailAddressResponse, error) {
	return s.service.ConfirmEmailAddress(ctx, request)
}
