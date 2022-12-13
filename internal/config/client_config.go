package config

import (
	"flag"
	"log"
	"os"
)

// ClientConfig represents client configuration
type ClientConfig struct {
	ServerAddress string
	UseTLS        bool
}

// GetClientConfig parses env variables and
// returns ClientConfig instance
func GetClientConfig(args []string) ClientConfig {
	var config ClientConfig
	var addressF string
	var useTLSf bool
	f := flag.NewFlagSet("server", flag.ExitOnError)
	f.StringVar(&addressF, "a", serverAddressDefault, "server address")
	f.BoolVar(&useTLSf, "t", false, "use TLS when connection to server")
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

	config.UseTLS = useTLSf

	return config
}
