package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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


	// parent context from which all request context will be derived
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// parsing config and preparing gRPC client
	clientConfig := config.GetClientConfig(os.Args[1:])
	mstorage := client.NewMemStorage()
	gclient, conn := client.NewGclient(clientConfig, &mstorage)
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	// Graceful shutown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s := <-sigCh
		log.Printf("got signal %v, attempting graceful shutdown", s)
		cancel()
		wg.Done()
	}()

	printStartupInfo()
	go client.Run(ctx, &gclient, &mstorage)
	wg.Wait()
	log.Println("Client was shutdown cleanelly")
}
