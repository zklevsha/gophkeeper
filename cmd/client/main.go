package main

import (
	"crypto/tls"
	"log"
	"os"
	"path"

	"github.com/zklevsha/gophkeeper/internal/client"
	"github.com/zklevsha/gophkeeper/internal/config"
	"github.com/zklevsha/gophkeeper/internal/structs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	// removing timestamps from the output
	log.SetFlags(0)

	mstorage := structs.MemStorage{}
	mstorage.MasterKeyDir = createKeyDir()

	clientConfig := config.GetClientConfig(os.Args[1:])
	// initiating connection to server
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	conn, err := grpc.Dial(clientConfig.ServerAddress,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithUnaryInterceptor(client.GetUnaryClientInterceptor(&mstorage)))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	gclient := structs.NewGclient(conn)
	// starting interactive loop
	client.Run(&gclient, &mstorage)
}

// createKeyDir creates MasterKeys directory (in needed) and returns its path
func createKeyDir() string {
	user_home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("cant get user home directory: %s", err.Error())
	}
	kpath := path.Join(user_home, ".gk-keychain")
	err = os.MkdirAll(kpath, 0700)
	if err != nil {
		log.Fatalf("cant create keychain directory(%s): %s", kpath, err.Error())
	}
	return kpath
}
