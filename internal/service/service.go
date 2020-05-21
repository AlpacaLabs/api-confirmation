package service

import (
	"github.com/AlpacaLabs/api-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-confirmation/internal/db"
	"google.golang.org/grpc"
)

type Service struct {
	config      configuration.Config
	dbClient    db.Client
	accountConn *grpc.ClientConn
}

func NewService(config configuration.Config, dbClient db.Client, accountConn *grpc.ClientConn) Service {
	return Service{
		config:      config,
		dbClient:    dbClient,
		accountConn: accountConn,
	}
}
