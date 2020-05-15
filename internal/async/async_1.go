package async

import (
	"context"
	"fmt"

	"github.com/AlpacaLabs/api-account-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-account-confirmation/internal/service"
	goKafka "github.com/AlpacaLabs/go-kafka"
	confirmationV1 "github.com/AlpacaLabs/protorepo-confirmation-go/alpacalabs/confirmation/v1"
	log "github.com/sirupsen/logrus"
)

const (
	TopicForEmailAddressConfirmation = "create-email-address-confirmation-code-request"
	TopicForPhoneNumberConfirmation  = "create-phone-number-confirmation-code-request"
)

func HandleCreateEmailAddressCode(config configuration.Config, s service.Service) {
	handle(TopicForEmailAddressConfirmation, config, handleCreateEmailAddressCode(s))
}

func HandleCreatePhoneNumberCode(config configuration.Config, s service.Service) {
	handle(TopicForPhoneNumberConfirmation, config, handleCreatePhoneNumberCode(s))
}

func handle(topic string, config configuration.Config, fn goKafka.ProcessFunc) {
	ctx := context.TODO()

	groupID := config.AppName
	brokers := []string{
		fmt.Sprintf("%s:%d", config.KafkaConfig.Host, config.KafkaConfig.Port),
	}

	err := goKafka.ProcessKafkaMessages(ctx, goKafka.ProcessKafkaMessagesInput{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		ProcessFunc: fn,
	})
	if err != nil {
		log.Errorf("%v", err)
	}
}

func handleCreateEmailAddressCode(s service.Service) goKafka.ProcessFunc {
	return func(ctx context.Context, message goKafka.Message) {
		// Convert kafka.Message to Protocol Buffer
		pb := &confirmationV1.CreateEmailAddressConfirmationCodeRequest{}
		if err := message.Unmarshal(pb); err != nil {
			log.Errorf("failed to unmarshal protobuf from kafka message: %v", err)
		}

		// TODO we could skip Kafka messages we've seen before by checking if
		//  the message's EventID exists in the TXOB table.

		if err := s.CreateEmailAddressConfirmationCode(ctx, pb); err != nil {
			log.Errorf("failed to process kafka message in transaction: %v", err)
		}
	}
}

func handleCreatePhoneNumberCode(s service.Service) goKafka.ProcessFunc {
	return func(ctx context.Context, message goKafka.Message) {
		// Convert kafka.Message to Protocol Buffer
		pb := &confirmationV1.CreatePhoneNumberConfirmationCodeRequest{}
		if err := message.Unmarshal(pb); err != nil {
			log.Errorf("failed to unmarshal protobuf from kafka message: %v", err)
		}

		// TODO we could skip Kafka messages we've seen before by checking if
		//  the message's EventID exists in the TXOB table.

		if err := s.CreatePhoneNumberConfirmationCode(ctx, pb); err != nil {
			log.Errorf("failed to process kafka message in transaction: %v", err)
		}
	}
}
