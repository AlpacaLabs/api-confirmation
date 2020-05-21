package async

import (
	"context"
	"fmt"
	"time"

	"github.com/AlpacaLabs/api-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-confirmation/internal/db"
	hermesTopics "github.com/AlpacaLabs/api-hermes/pkg/topic"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

func RelayMessagesForSendEmail(config configuration.Config, dbClient db.Client) {
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
		fn := relayMessage(relayMessageInput{
			topic:                    topic,
			writer:                   writer,
			transactionalOutboxTable: db.TableForSendEmailRequest,
		})
		err := dbClient.RunInTransaction(ctx, fn)
		if err != nil {
			log.Errorf("message relay encountered error... sleeping for a bit... %v", err)
			time.Sleep(time.Second * 2)
		}
	}
}
