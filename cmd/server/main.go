package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/zklevsha/gophkeeper/internal/config"
	"github.com/zklevsha/gophkeeper/internal/db"
	"github.com/zklevsha/gophkeeper/internal/gserver"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func printStartupInfo() {
	log.Printf("build version: %s, build date: %s, build commit: %s",
		buildVersion, buildDate, buildCommit)
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second*60))
	defer cancel()

	// initialize db.Connector and gRPC servser instance
	config := config.GetServerConfig(os.Args[1:])
	db := db.Connector{DSN: config.DSN,}
	err := db.Init(ctx)
	if err != nil {
		log.Fatalf("failed to initialize db.Connector: %s", err.Error())
	}
	gserver := gserver.GetServer(db, config.Key, config.CertPath,
		config.PrivateKeyPath)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	wg := sync.WaitGroup{}

	// signal listener
	wg.Add(1)
	go func() {
		s := <-sigCh
		log.Printf("got signal %v, attempting graceful shutdown", s)
		cancel()
		gserver.GracefulStop()
		// grpc.Stop() // leads to error while receiving stream response: rpc error: code = Unavailable desc = transport is closing
		wg.Done()
	}()

	// stating gRPC server
	log.Printf("starting server on %s", config.ServerAddress)
	printStartupInfo()
	listener, err := net.Listen("tcp", config.ServerAddress)
	if err != nil {
		panic(err)
	}
	if err := gserver.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	wg.Wait()
	log.Println("clean shutdown")
}
