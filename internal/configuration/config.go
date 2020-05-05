package configuration

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	flag "github.com/spf13/pflag"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	flagForDBUser         = "db_user"
	flagForDBPass         = "db_pass"
	flagForDBHost         = "db_host"
	flagForDBName         = "db_name"
	flagForGrpcPort       = "grpc_port"
	flagForGrpcPortHealth = "grpc_port_health"
	flagForHTTPPort       = "http_port"

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

	DBHost string
	DBUser string
	DBPass string
	DBName string

	// GrpcPort controls what port our gRPC server runs on.
	GrpcPort int

	// HTTPPort controls what port our HTTP server runs on.
	HTTPPort int

	// HealthPort controls what port our gRPC health endpoints run on.
	HealthPort int

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
		AppName:    "citadel",
		AppID:      uuid.New().String(),
		GrpcPort:   8081,
		HealthPort: 8082,
		HTTPPort:   8083,
	}

	flag.String(flagForDBUser, c.DBUser, "DB user")
	flag.String(flagForDBPass, c.DBPass, "DB pass")
	flag.String(flagForDBHost, c.DBHost, "DB host")
	flag.String(flagForDBName, c.DBName, "DB name")

	flag.Int(flagForGrpcPort, c.GrpcPort, "gRPC port")
	flag.Int(flagForGrpcPortHealth, c.HealthPort, "gRPC health port")
	flag.Int(flagForHTTPPort, c.HTTPPort, "gRPC HTTP port")

	flag.String(flagForAccountGrpcAddress, "", "Address of Account gRPC service")
	flag.String(flagForAccountGrpcHost, "", "Host of Account gRPC service")
	flag.String(flagForAccountGrpcPort, "", "Port of Account gRPC service")

	flag.String(flagForHermesGrpcAddress, "", "Address of Hermes gRPC service")
	flag.String(flagForHermesGrpcHost, "", "Host of Hermes gRPC service")
	flag.String(flagForHermesGrpcPort, "", "Port of Hermes gRPC service")

	flag.Parse()

	viper.BindPFlag(flagForDBUser, flag.Lookup(flagForDBUser))
	viper.BindPFlag(flagForDBPass, flag.Lookup(flagForDBPass))
	viper.BindPFlag(flagForDBHost, flag.Lookup(flagForDBHost))
	viper.BindPFlag(flagForDBName, flag.Lookup(flagForDBName))

	viper.BindPFlag(flagForGrpcPort, flag.Lookup(flagForGrpcPort))
	viper.BindPFlag(flagForGrpcPortHealth, flag.Lookup(flagForGrpcPortHealth))
	viper.BindPFlag(flagForHTTPPort, flag.Lookup(flagForHTTPPort))

	viper.BindPFlag(flagForAccountGrpcAddress, flag.Lookup(flagForAccountGrpcAddress))
	viper.BindPFlag(flagForAccountGrpcHost, flag.Lookup(flagForAccountGrpcHost))
	viper.BindPFlag(flagForAccountGrpcPort, flag.Lookup(flagForAccountGrpcPort))

	viper.BindPFlag(flagForHermesGrpcAddress, flag.Lookup(flagForHermesGrpcAddress))
	viper.BindPFlag(flagForHermesGrpcHost, flag.Lookup(flagForHermesGrpcHost))
	viper.BindPFlag(flagForHermesGrpcPort, flag.Lookup(flagForHermesGrpcPort))

	viper.AutomaticEnv()

	c.DBUser = viper.GetString(flagForDBUser)
	c.DBPass = viper.GetString(flagForDBPass)
	c.DBHost = viper.GetString(flagForDBHost)
	c.DBName = viper.GetString(flagForDBName)

	c.GrpcPort = viper.GetInt(flagForGrpcPort)
	c.HealthPort = viper.GetInt(flagForGrpcPortHealth)
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
