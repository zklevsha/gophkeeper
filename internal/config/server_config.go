package config

import (
	"flag"
	"os"
)

type ServerConfig struct {
	ServerAddress string
	DSN           string
}

const serverAddressDefault = "localhost:443"

func GetServerConfig(args []string) ServerConfig {
	var config ServerConfig
	var addressF, DSNf string
	f := flag.NewFlagSet("server", flag.ExitOnError)
	f.StringVar(&addressF, "a", serverAddressDefault, "server address")
	f.StringVar(&DSNf, "d", "", "database connection string (postgres://username:password@localhost:5432/database_name)")
	f.Parse(args)

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
