package async

import (
	"github.com/AlpacaLabs/api-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-confirmation/internal/db"
	hermesTopics "github.com/AlpacaLabs/api-hermes/pkg/topic"
)

func RelayMessagesForSendEmail(config configuration.Config, dbClient db.Client) {
	relayMessages(config, dbClient, relayMessagesInput{
		topic:                    hermesTopics.TopicForSendEmailRequest,
		transactionalOutboxTable: db.TableForSendEmailRequest,
	})
}

func RelayMessagesForSendSms(config configuration.Config, dbClient db.Client) {
	relayMessages(config, dbClient, relayMessagesInput{
		topic:                    hermesTopics.TopicForSendSmsRequest,
		transactionalOutboxTable: db.TableForSendSmsRequest,
	})
}

func RelayMessagesForConfirmEmail(config configuration.Config, dbClient db.Client) {
	relayMessages(config, dbClient, relayMessagesInput{
		// TODO get topic name from Account service
		topic:                    "confirm-email-request",
		transactionalOutboxTable: db.TableForConfirmEmailAddressRequest,
	})
}

func RelayMessagesForConfirmPhone(config configuration.Config, dbClient db.Client) {
	relayMessages(config, dbClient, relayMessagesInput{
		// TODO get topic name from Account service
		topic:                    "confirm-phone-request",
		transactionalOutboxTable: db.TableForConfirmPhoneNumberRequest,
	})
}
