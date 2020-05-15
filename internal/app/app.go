package app

import (
	"sync"

	"github.com/AlpacaLabs/api-account-confirmation/internal/grpc"

	"github.com/AlpacaLabs/api-account-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-account-confirmation/internal/db"
	"github.com/AlpacaLabs/api-account-confirmation/internal/http"
	"github.com/AlpacaLabs/api-account-confirmation/internal/service"
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
	dbConn, err := db.Connect(a.config.SQLConfig)
	if err != nil {
		log.Fatalf("failed to dial database: %v", err)
	}
	dbClient := db.NewClient(dbConn)
	svc := service.NewService(a.config, dbClient)

	var wg sync.WaitGroup

	wg.Add(1)
	httpServer := http.NewServer(a.config, svc)
	httpServer.Run()

	wg.Add(1)
	grpcServer := grpc.NewServer(a.config, svc)
	grpcServer.Run()

	wg.Wait()
}
