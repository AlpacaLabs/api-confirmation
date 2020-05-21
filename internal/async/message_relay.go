package async

import (
	"context"
	"fmt"

	"github.com/AlpacaLabs/api-confirmation/internal/db"
	goKafka "github.com/AlpacaLabs/go-kafka"
	"github.com/segmentio/kafka-go"
)

type relayMessageInput struct {
	topic                    string
	writer                   *kafka.Writer
	transactionalOutboxTable string
}

func relayMessages() {
	// TODO DRY this out
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
