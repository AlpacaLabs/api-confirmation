package async

import (
	"context"
	"fmt"

	"github.com/AlpacaLabs/api-account-confirmation/internal/configuration"
	goKafka "github.com/AlpacaLabs/go-kafka"
	log "github.com/sirupsen/logrus"
)

func readFromTopic(topic string, config configuration.Config, fn goKafka.ProcessFunc) {
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
