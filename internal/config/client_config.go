package config

import (
	"flag"
	"log"
	"os"
)

// ClientConfig represents client configuration
type ClientConfig struct {
	ServerAddress string
}

// GetClientConfig parses env variables and
// returns ClientConfig instance
func GetClientConfig(args []string) ClientConfig {
	var config ClientConfig
	var addressF string
	f := flag.NewFlagSet("server", flag.ExitOnError)
	f.StringVar(&addressF, "a", serverAddressDefault, "server address")
	err := f.Parse(args)
	if err != nil {
		log.Fatalf("failed to parse flag parameters: %s",
			err.Error())
	}

	addressEnv := os.Getenv("GK_SERVER_ADDRESS")
	if addressEnv != "" {
		config.ServerAddress = addressEnv
	} else {
		config.ServerAddress = addressF
	}

	return config
}
