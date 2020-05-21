package app

import (
	"sync"

	"github.com/AlpacaLabs/go-kontext"

	"github.com/AlpacaLabs/api-confirmation/internal/async"

	"github.com/AlpacaLabs/api-confirmation/internal/grpc"

	"github.com/AlpacaLabs/api-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-confirmation/internal/db"
	"github.com/AlpacaLabs/api-confirmation/internal/http"
	"github.com/AlpacaLabs/api-confirmation/internal/service"
	log "github.com/sirupsen/logrus"
)

type App struct {
	config configuration.Config
}

func NewApp(c configuration.Config) App {
	return App{
		config: c,
	}
}

func (a App) Run() {
	config := a.config

	// Connect to database
	dbConn, err := db.Connect(config.SQLConfig)
	if err != nil {
		log.Fatalf("failed to dial database: %v", err)
	}
	dbClient := db.NewClient(dbConn)

	// Connect to Account service
	accountConn, err := kontext.Dial(config.AccountGRPCAddress)
	if err != nil {
		log.Fatalf("failed to dial Account service: %v", err)
	}

	// Create our service layer
	svc := service.NewService(a.config, dbClient, accountConn)

	var wg sync.WaitGroup

	wg.Add(1)
	httpServer := http.NewServer(a.config, svc)
	go httpServer.Run()

	wg.Add(1)
	grpcServer := grpc.NewServer(a.config, svc)
	go grpcServer.Run()

	wg.Add(1)
	go async.ConsumeTopicForEmailCodeCreation(a.config, svc)

	wg.Add(1)
	go async.ConsumeTopicForPhoneCodeCreation(a.config, svc)

	wg.Add(1)
	go async.RelayMessagesForSendEmail(a.config, dbClient)

	wg.Add(1)
	go async.RelayMessagesForSendSms(a.config, dbClient)

	wg.Add(1)
	go async.RelayMessagesForConfirmEmail(a.config, dbClient)

	wg.Add(1)
	go async.RelayMessagesForConfirmPhone(a.config, dbClient)

	wg.Wait()
}
