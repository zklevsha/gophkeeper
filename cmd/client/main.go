package main

import (
	"log"
	"os"

	"github.com/zklevsha/gophkeeper/internal/client"
	"github.com/zklevsha/gophkeeper/internal/config"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func printStartupInfo() {
	log.Printf("build version: %s, build date: %s, build commit: %s",
		buildVersion, buildDate, buildCommit)
}


func main() {
	// removing timestamps from the output
	log.SetFlags(0)

	clientConfig := config.GetClientConfig(os.Args[1:])
	mstorage := client.NewMemStorage()
	gclient, conn := client.NewGclient(clientConfig, &mstorage)

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// starting interactive loop
	printStartupInfo()
	client.Run(&gclient, &mstorage)
}
