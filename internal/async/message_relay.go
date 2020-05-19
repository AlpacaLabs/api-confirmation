package async

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	hermesTopics "github.com/AlpacaLabs/api-hermes/pkg/topic"

	"github.com/AlpacaLabs/api-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-confirmation/internal/db"
	"github.com/segmentio/kafka-go"
	log "github.com/sirupsen/logrus"
)

type relayMessageInput struct {
	topic  string
	writer *kafka.Writer
}

func RelayMessagesForSendEmail(config configuration.Config, dbClient db.Client, accountConn *grpc.ClientConn) {
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
		fn := relayMessageForSendEmail(relayMessageForSendEmailInput{
			accountConn: accountConn,
			relayMessageInput: relayMessageInput{
				topic:  topic,
				writer: writer,
			},
		})
		err := dbClient.RunInTransaction(ctx, fn)
		if err != nil {
			log.Errorf("message relay encountered error... sleeping for a bit... %v", err)
			time.Sleep(time.Second * 2)
		}
	}
}

func RelayMessagesForSendSms(config configuration.Config, dbClient db.Client, accountConn *grpc.ClientConn) {
	topic := hermesTopics.TopicForSendSmsRequest
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
		fn := relayMessageForSendSms(relayMessageForSendSmsInput{
			accountConn: accountConn,
			relayMessageInput: relayMessageInput{
				topic:  topic,
				writer: writer,
			},
		})
		err := dbClient.RunInTransaction(ctx, fn)
		if err != nil {
			log.Errorf("message relay encountered error... sleeping for a bit... %v", err)
			time.Sleep(time.Second * 2)
		}
	}
}
