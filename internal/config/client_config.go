package config

import "os"

type ClientConfig struct {
	ServerAddress string `json:"address,omitempty"`
}

// GetClientConfig rep
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
