package async

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	hermesTopics "github.com/AlpacaLabs/api-hermes/pkg/topic"

	"github.com/AlpacaLabs/api-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-confirmation/internal/db"
	goKafka "github.com/AlpacaLabs/go-kafka"
	accountV1 "github.com/AlpacaLabs/protorepo-account-go/alpacalabs/account/v1"
	hermesV1 "github.com/AlpacaLabs/protorepo-hermes-go/alpacalabs/hermes/v1"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

// As described in the docs: https://microservices.io/patterns/data/transactional-outbox.html
// "A separate Message Relay process publishes the events inserted into database to a message broker."

func ReadFromTransactionalOutbox(config configuration.Config, dbClient db.Client, accountConn *grpc.ClientConn) {
	topic := hermesTopics.TopicForSendEmailRequest
	brokers := []string{
		fmt.Sprintf("%s:%d", config.KafkaConfig.Host, config.KafkaConfig.Port),
	}
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic,
	})
	defer writer.Close()

	for {
		ctx := context.TODO()
		err := dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
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
		})
		if err != nil {
			log.Errorf("message relay encountered error... sleeping for a bit... %v", err)
			time.Sleep(time.Second * 2)
		}
	}
}
