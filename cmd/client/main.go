package main

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/zklevsha/gophkeeper/internal/client"
	"github.com/zklevsha/gophkeeper/internal/config"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// removing timestamps from the output
	log.SetFlags(0)

	clientConfig := config.GetClientConfig(os.Args[1:])

	// initiating connection to server
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := grpc.Dial(clientConfig.ServerAddress, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	gclient := structs.Gclient{Auth: pb.NewAuthClient(conn)}

	// starting interactive loop
	client.Run(gclient)
}
