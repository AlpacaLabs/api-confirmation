package async

import (
	"context"
	"fmt"

	"github.com/AlpacaLabs/api-confirmation/internal/db"
	goKafka "github.com/AlpacaLabs/go-kafka"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	hermesV1 "github.com/AlpacaLabs/protorepo-hermes-go/alpacalabs/hermes/v1"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type relayMessageForSendSmsInput struct {
	accountConn *grpc.ClientConn
	relayMessageInput
}

func relayMessageForSendSms(in relayMessageForSendSmsInput) db.TransactionFunc {
	accountConn := in.accountConn
	writer := in.writer
	topic := in.topic

	return func(ctx context.Context, tx db.Transaction) error {
		e, err := tx.ReadEventForSendSmsRequest(ctx)
		if err != nil {
			return fmt.Errorf("failed to read event from transactional outbox for sending emails: %w", err)
		}

		codeEntity, err := tx.GetPhoneCodeByID(ctx, e.CodeID)
		if err != nil {
			return err
		}

		log.Infof("Looking up email address for email code: %s", codeEntity.Id)

		var phoneNumber string
		accountClient := accountV1.NewAccountServiceClient(accountConn)
		if res, err := accountClient.GetPhoneNumber(ctx, &accountV1.GetPhoneNumberRequest{
			Id: codeEntity.PhoneNumberId,
		}); err != nil {
			return fmt.Errorf("failed to verify email address existence: %w", err)
		} else {
			phoneNumber = res.PhoneNumber.PhoneNumber
		}

		pb := &hermesV1.SendSmsRequest{
			To:      phoneNumber,
			Message: fmt.Sprintf("Your confirmation code for this phone number is: %s", codeEntity.Code),
		}

		msg, err := goKafka.NewMessage(e.TraceInfo, e.EventInfo, pb)
		if err != nil {
			return fmt.Errorf("failed to create event for topic: %s: %w", topic, err)
		}

		if err := writer.WriteMessages(ctx, msg.ToKafkaMessage()); err != nil {
			return fmt.Errorf("failed to send error to topic: %s: %w", topic, err)
		}

		return tx.MarkAsSentSendSmsRequest(ctx, e.EventId)
	}
}
