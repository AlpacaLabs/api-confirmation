package service

import (
	"github.com/AlpacaLabs/api-account-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-account-confirmation/internal/db"
)

type Service struct {
	config   configuration.Config
	dbClient db.Client
}

func NewService(config configuration.Config, dbClient db.Client) Service {
	return Service{
		config:   config,
		dbClient: dbClient,
	}
}
