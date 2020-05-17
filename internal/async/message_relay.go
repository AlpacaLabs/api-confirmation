package async

import (
	"context"
	"fmt"
	"time"

	"github.com/AlpacaLabs/api-hermes/pkg/topic"

	"github.com/AlpacaLabs/api-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-confirmation/internal/db"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

// As described in the docs: https://microservices.io/patterns/data/transactional-outbox.html
// "A separate Message Relay process publishes the events inserted into database to a message broker."

func ReadFromTransactionalOutbox(config configuration.Config, dbClient db.Client) {
	brokers := []string{
		fmt.Sprintf("%s:%d", config.KafkaConfig.Host, config.KafkaConfig.Port),
	}
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: brokers,
		Topic:   topic.TopicForSendEmailRequest,
	})
	defer writer.Close()

	for {
		ctx := context.TODO()
		if err := dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
			e, err := tx.ReadEventForSendEmailRequest(ctx)
			if err != nil {
				return fmt.Errorf("failed to read event from transactional outbox for sending emails: %w", err)
			}

			// TODO write protobuf to Hermes topic
			//writer.WriteMessages(ctx,)

			return tx.MarkAsSentSendEmailRequest(ctx, e.EventId)
		}); err != nil {
			log.Errorf("message relay encountered error: %v", err)
			log.Warnf("sleeping after error...")
			time.Sleep(time.Second * 2)
			break
		}
	}
}
