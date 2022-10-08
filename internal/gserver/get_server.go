package gserver

import (
	"crypto/tls"
	"log"

	"github.com/zklevsha/gophkeeper/internal/db"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

func loadTLSCredentials(certPath string, privatePath string) (credentials.TransportCredentials, error) {
	// Load server's certificate and private key
	serverCert, err := tls.LoadX509KeyPair(certPath, privatePath)
	if err != nil {
		return nil, err
	}

	// Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

// GetServer registers all gRPC services and returns grpc.Server instanse
func GetServer(db db.Connector, key string, certPath string, privatePath string) *grpc.Server {
	var s *grpc.Server
	if certPath != "" && privatePath != "" {
		tlsCredentials, err := loadTLSCredentials(certPath, privatePath)
		if err != nil {
			log.Fatal("cannot load TLS credentials: ", err)
		}

		s = grpc.NewServer(
			grpc.Creds(tlsCredentials),
		)
	} else {
		s = grpc.NewServer()
	}

	reflection.Register(s)
	pb.RegisterAuthServer(s, &authServer{db: db, key: key})
	return s

}
