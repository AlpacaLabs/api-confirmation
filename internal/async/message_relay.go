package async

import (
	"context"
	"fmt"
	"time"

	"github.com/AlpacaLabs/api-hermes/pkg/topic"

	"github.com/AlpacaLabs/api-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-confirmation/internal/db"
	"github.com/AlpacaLabs/api-confirmation/internal/db/entities"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

// As described in the docs: https://microservices.io/patterns/data/transactional-outbox.html
// "A separate Message Relay process publishes the events inserted into database to a message broker."

func ReadFromTransactionalOutbox(config configuration.Config, dbClient db.Client) {
	brokers := []string{
		fmt.Sprintf("%s:%d", config.KafkaConfig.Host, config.KafkaConfig.Port),
	}
	batchSize := 5
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:   brokers,
		Topic:     topic.TopicForSendEmailRequest,
		BatchSize: batchSize,
	})
	defer writer.Close()
	for {
		ctx := context.TODO()
		var events []entities.SendEmailEvent
		if err := dbClient.RunInTransaction(ctx, func(ctx context.Context, tx db.Transaction) error {
			e, err := tx.ReadFromTxobForEmailCode(ctx, batchSize)
			events = e
			return err
		}); err != nil {
			log.Errorf("message relay encountered error: %v", err)
			log.Warnf("sleeping after error...")
			time.Sleep(time.Second * 2)
			break
		}

		log.Debugf("Sending events: %v", events)

		//for _, e := range events {
		// TODO send to topic
		//writer.WriteMessages(ctx, )
	}
}
