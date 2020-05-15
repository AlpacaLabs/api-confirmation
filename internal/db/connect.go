package db

import (
	"context"
	"fmt"

	configuration "github.com/AlpacaLabs/go-config"

	"github.com/jackc/pgx/v4"
)

func Connect(config configuration.SQLConfig) (*pgx.Conn, error) {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable",
		config.User, config.Pass, config.Host, config.Name)

	return pgx.Connect(context.TODO(), connectionString)
}
