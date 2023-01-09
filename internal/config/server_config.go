package config

import (
	"flag"
	"log"
	"os"
)

// ServerConfig represents server configuration
type ServerConfig struct {
	ServerAddress  string
	DSN            string
	Key            string
	CertPath       string
	PrivateKeyPath string
}

const serverAddressDefault = "localhost:443"

// GetServerConfig parses runtime flags and
// enviroment variables and returns ServierConfig instanse
func GetServerConfig(args []string) ServerConfig {
	var config ServerConfig
	var addressF, DSNf, keyF, certPathF, privateKeyPathF string
	f := flag.NewFlagSet("server", flag.ExitOnError)
	f.StringVar(&addressF, "a", serverAddressDefault, "server address")
	f.StringVar(&DSNf, "d", "", "database connection string (postgres://username:password@localhost:5432/database_name)")
	f.StringVar(&keyF, "k", "", "server key to sign JWT tokens with")
	f.StringVar(&certPathF, "c", "", "path to server`s certificate (if not set data will not be encrypted)")
	f.StringVar(&privateKeyPathF, "p", "", "path to server`s private key (if not set data will not be encrypted)")
	err := f.Parse(args)
	if err != nil {
		log.Fatalf("failed to parse flag parameters: %s",
			err.Error())
	}

	addressEnv := os.Getenv("GK_SERVER_ADDRESS")
	DSNenv := os.Getenv("GK_DB_DSN")
	keyEnv := os.Getenv("GK_KEY")
	certPathEnv := os.Getenv("GK_CERT")
	privateKeyPathEnv := os.Getenv("GK_PRIVATE_KEY")

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
	if config.DSN == "" {
		log.Fatalf("DSN is not set")
	}

	// Key
	if keyEnv != "" {
		config.Key = keyEnv
	} else {
		config.Key = keyF
	}
	if config.Key == "" {
		log.Fatal("KEY is not set")
	}

	// CertPath
	if certPathEnv != "" {
		config.CertPath = certPathEnv
	} else {
		config.CertPath = certPathF
	}

	// PrivatePath
	if privateKeyPathEnv != "" {
		config.PrivateKeyPath = privateKeyPathEnv
	} else {
		config.PrivateKeyPath = privateKeyPathF
	}

	return config
}
