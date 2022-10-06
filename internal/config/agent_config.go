package config

import "os"

type AgentConfig struct {
	ServerAddress string `json:"address,omitempty"`
}

func GetAgentConfig() AgentConfig {
	var config AgentConfig

	addressEnv := os.Getenv("GK_SERVER_ADDRESS")
	if addressEnv != "" {
		config.ServerAddress = addressEnv
	} else {
		config.ServerAddress = serverAddressDefault
	}

	return config
}
