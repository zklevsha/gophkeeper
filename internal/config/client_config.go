package config

import "os"

// ClientConfig represents client configuration
type ClientConfig struct {
	ServerAddress string
}

// GetClientConfig parses env variables and
// returns ClientConfig instance
func GetClientConfig() ClientConfig {
	var config ClientConfig

	addressEnv := os.Getenv("GK_SERVER_ADDRESS")
	if addressEnv != "" {
		config.ServerAddress = addressEnv
	} else {
		config.ServerAddress = serverAddressDefault
	}

	return config
}
