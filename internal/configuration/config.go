package configuration

import (
	"encoding/json"
	"fmt"

	configuration "github.com/AlpacaLabs/go-config"
	"github.com/rs/xid"

	flag "github.com/spf13/pflag"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	flagForGrpcPort = "grpc_port"
	flagForHTTPPort = "http_port"

	flagForAccountGrpcAddress = "account_service_address"
	flagForAccountGrpcHost    = "account_service_host"
	flagForAccountGrpcPort    = "account_service_port_grpc"

	flagForHermesGrpcAddress = "hermes_service_address"
	flagForHermesGrpcHost    = "hermes_service_host"
	flagForHermesGrpcPort    = "hermes_service_port_grpc"
)

type Config struct {
	// AppName is a low cardinality identifier for this service.
	AppName string

	// AppID is a unique identifier for the instance (pod) running this app.
	AppID string

	// KafkaConfig provides configuration for connecting to Apache Kafka.
	KafkaConfig configuration.KafkaConfig

	// SQLConfig provides configuration for connecting to a SQL database.
	SQLConfig configuration.SQLConfig

	// GrpcPort controls what port our gRPC server runs on.
	GrpcPort int

	// HTTPPort controls what port our HTTP server runs on.
	HTTPPort int

	// AccountGRPCAddress is the gRPC address of the Account service.
	AccountGRPCAddress string

	// HermesGRPCAddress is the gRPC address of the Hermes service.
	HermesGRPCAddress string
}

func (c Config) String() string {
	b, err := json.Marshal(c)
	if err != nil {
		log.Fatalf("Could not marshal config to string: %v", err)
	}
	return string(b)
}

func LoadConfig() Config {
	c := Config{
		AppName:  "api-confirmation",
		AppID:    xid.New().String(),
		GrpcPort: 8081,
		HTTPPort: 8083,
	}

	c.KafkaConfig = configuration.LoadKafkaConfig()
	c.SQLConfig = configuration.LoadSQLConfig()

	flag.Int(flagForGrpcPort, c.GrpcPort, "gRPC port")
	flag.Int(flagForHTTPPort, c.HTTPPort, "HTTP port")

	flag.String(flagForAccountGrpcAddress, "", "Address of Account gRPC service")
	flag.String(flagForAccountGrpcHost, "", "Host of Account gRPC service")
	flag.String(flagForAccountGrpcPort, "", "Port of Account gRPC service")

	flag.String(flagForHermesGrpcAddress, "", "Address of Hermes gRPC service")
	flag.String(flagForHermesGrpcHost, "", "Host of Hermes gRPC service")
	flag.String(flagForHermesGrpcPort, "", "Port of Hermes gRPC service")

	flag.Parse()

	viper.BindPFlag(flagForGrpcPort, flag.Lookup(flagForGrpcPort))
	viper.BindPFlag(flagForHTTPPort, flag.Lookup(flagForHTTPPort))

	viper.BindPFlag(flagForAccountGrpcAddress, flag.Lookup(flagForAccountGrpcAddress))
	viper.BindPFlag(flagForAccountGrpcHost, flag.Lookup(flagForAccountGrpcHost))
	viper.BindPFlag(flagForAccountGrpcPort, flag.Lookup(flagForAccountGrpcPort))

	viper.BindPFlag(flagForHermesGrpcAddress, flag.Lookup(flagForHermesGrpcAddress))
	viper.BindPFlag(flagForHermesGrpcHost, flag.Lookup(flagForHermesGrpcHost))
	viper.BindPFlag(flagForHermesGrpcPort, flag.Lookup(flagForHermesGrpcPort))

	viper.AutomaticEnv()

	c.GrpcPort = viper.GetInt(flagForGrpcPort)
	c.HTTPPort = viper.GetInt(flagForHTTPPort)

	c.AccountGRPCAddress = getGrpcAddress(flagForAccountGrpcAddress, flagForAccountGrpcHost, flagForAccountGrpcPort)
	c.HermesGRPCAddress = getGrpcAddress(flagForHermesGrpcAddress, flagForHermesGrpcHost, flagForHermesGrpcPort)

	return c
}

func getGrpcAddress(addrFlag, hostFlag, portFlag string) string {
	addr := viper.GetString(addrFlag)
	host := viper.GetString(hostFlag)
	port := viper.GetInt(portFlag)

	if port != 0 {
		return fmt.Sprintf("%s:%d", host, port)
	}

	return addr
}
