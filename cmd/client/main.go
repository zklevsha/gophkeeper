package main

import (
	"log"
	"os"

	"github.com/zklevsha/gophkeeper/internal/client"
	"github.com/zklevsha/gophkeeper/internal/config"
)

func main() {
	// removing timestamps from the output
	log.SetFlags(0)

	clientConfig := config.GetClientConfig(os.Args[1:])
	mstorage := client.NewMemStorage()
	gclient := client.NewGclient(clientConfig, mstorage)

	// starting interactive loop
	client.Run(&gclient, &mstorage)
}
