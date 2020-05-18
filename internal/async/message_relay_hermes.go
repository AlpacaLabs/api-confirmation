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

type relayMessageToHermesInput struct {
	accountConn *grpc.ClientConn
	relayMessageInput
}

func relayMessageToHermes(in relayMessageToHermesInput) db.TransactionFunc {
	accountConn := in.accountConn
	writer := in.writer
	topic := in.topic

	return func(ctx context.Context, tx db.Transaction) error {
		e, err := tx.ReadEventForSendEmailRequest(ctx)
		if err != nil {
			return fmt.Errorf("failed to read event from transactional outbox for sending emails: %w", err)
		}

		emailCode, err := tx.GetEmailCodeByID(ctx, e.CodeID)
		if err != nil {
			return err
		}

		log.Infof("Looking up email address for email code: %s", emailCode.Id)

		var emailAddress string
		accountClient := accountV1.NewAccountServiceClient(accountConn)
		if emailRes, err := accountClient.GetEmailAddress(ctx, &accountV1.GetEmailAddressRequest{
			Id: emailCode.EmailAddressId,
		}); err != nil {
			return fmt.Errorf("failed to verify email address existence: %w", err)
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

		msg, err := goKafka.NewMessage(e.TraceInfo, e.EventInfo, pb)
		if err != nil {
			return fmt.Errorf("failed to create event for topic: %s: %w", topic, err)
		}

		if err := writer.WriteMessages(ctx, msg.ToKafkaMessage()); err != nil {
			return fmt.Errorf("failed to send error to topic: %s: %w", topic, err)
		}

		return tx.MarkAsSentSendEmailRequest(ctx, e.EventId)
	}
}
