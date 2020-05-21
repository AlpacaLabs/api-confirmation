package async

import (
	"context"
	"fmt"
	"time"

	"github.com/AlpacaLabs/api-confirmation/internal/configuration"
	log "github.com/sirupsen/logrus"

	"github.com/AlpacaLabs/api-confirmation/internal/db"
	goKafka "github.com/AlpacaLabs/go-kafka"
	"github.com/segmentio/kafka-go"
)

type relayMessagesInput struct {
	topic                    string
	transactionalOutboxTable string
}

func relayMessages(config configuration.Config, dbClient db.Client, in relayMessagesInput) {
	topic := in.topic
	transactionalOutboxTable := in.transactionalOutboxTable

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
		fn := relayMessage(relayMessageInput{
			topic:                    topic,
			writer:                   writer,
			transactionalOutboxTable: transactionalOutboxTable,
		})
		err := dbClient.RunInTransaction(ctx, fn)
		if err != nil {
			log.Errorf("message relay encountered error... sleeping for a bit... %v", err)
			time.Sleep(time.Second * 2)
		}
	}
}

type relayMessageInput struct {
	topic                    string
	writer                   *kafka.Writer
	transactionalOutboxTable string
}

func relayMessage(in relayMessageInput) db.TransactionFunc {
	writer := in.writer
	topic := in.topic
	transactionalOutboxTable := in.transactionalOutboxTable

	return func(ctx context.Context, tx db.Transaction) error {
		e, err := tx.ReadEvent(ctx, transactionalOutboxTable)
		if err != nil {
			return fmt.Errorf("failed to read event from transactional outbox '%s': %w", transactionalOutboxTable, err)
		}

		msg, err := goKafka.NewMessage(e.TraceInfo, e.EventInfo, e.Payload)
		if err != nil {
			return fmt.Errorf("failed to create event for topic: %s: %w", topic, err)
		}

		if err := writer.WriteMessages(ctx, msg.ToKafkaMessage()); err != nil {
			return fmt.Errorf("failed to send error to topic: %s: %w", topic, err)
		}

		return tx.MarkEventAsSent(ctx, e.EventId, transactionalOutboxTable)
	}
}
