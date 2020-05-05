package db

import (
	"database/sql"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func Connect(user, password, host, dbname string) *sql.DB {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s dbname=%s sslmode=disable", user, password, host, dbname)

	dbConn, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}

	return dbConn
}
