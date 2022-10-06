package config

import (
	"flag"
	"log"
	"os"
)

// ServerConfig represents server configuration
type ServerConfig struct {
	ServerAddress string
	DSN           string
}

const serverAddressDefault = "localhost:443"

// GetServerConfig parses runtime flags and
// enviroment variables and returns ServierConfig instanse
func GetServerConfig(args []string) ServerConfig {
	var config ServerConfig
	var addressF, DSNf string
	f := flag.NewFlagSet("server", flag.ExitOnError)
	f.StringVar(&addressF, "a", serverAddressDefault, "server address")
	f.StringVar(&DSNf, "d", "", "database connection string (postgres://username:password@localhost:5432/database_name)")
	err := f.Parse(args)
	if err != nil {
		log.Fatalf("failed to parse flag parameters: %s",
			err.Error())
	}

	addressEnv := os.Getenv("GK_SERVER_ADDRESS")
	DSNenv := os.Getenv("GK_DB_DSN")

	// address
	if addressEnv != "" {
		config.ServerAddress = addressEnv
	} else {
		config.ServerAddress = addressF
	}

	// DSN
	if DSNenv != "" {
		config.DSN = DSNenv
	} else {
		config.DSN = DSNf
	}
	return config
}
