package app

import (
	"sync"

	"github.com/AlpacaLabs/api-account-confirmation/internal/grpc"

	"github.com/AlpacaLabs/api-account-confirmation/internal/configuration"
	"github.com/AlpacaLabs/api-account-confirmation/internal/db"
	"github.com/AlpacaLabs/api-account-confirmation/internal/http"
	"github.com/AlpacaLabs/api-account-confirmation/internal/service"
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
	dbConn := db.Connect(a.config.DBUser, a.config.DBPass, a.config.DBHost, a.config.DBName)
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
