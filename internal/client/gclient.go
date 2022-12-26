package client

import (
	"context"
	"crypto/tls"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/zklevsha/gophkeeper/internal/config"
	"github.com/zklevsha/gophkeeper/internal/pb"
)

// Gclient represents collection of various gRPC clients
type Gclient struct {
	Auth  pb.AuthClient
	Pdata pb.PrivateDataClient
}

// NewGclient initializes new Gclient
func NewGclient(clientConfig config.ClientConfig, mstorage MemStorage) Gclient {


	var conn *grpc.ClientConn
	var err error
	if clientConfig.UseTLS {
		tlsConfig := &tls.Config{
			InsecureSkipVerify: true,
		}
		conn, err = grpc.Dial(clientConfig.ServerAddress,
			grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
			grpc.WithUnaryInterceptor(getUnaryClientInterceptor(&mstorage)))
	} else {
		conn, err = grpc.Dial(clientConfig.ServerAddress,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(getUnaryClientInterceptor(&mstorage)))

	}
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer func() {
		err := conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()


	return Gclient{
		Auth:  pb.NewAuthClient(conn),
		Pdata: pb.NewPrivateDataClient(conn)}
}


// list of RPC that not required authorization
var noAuth = map[string]bool{
	"/Auth/Register": true,
	"/Auth/GetToken": true,
}

func getUnaryClientInterceptor(mstorage *MemStorage) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if !noAuth[method] {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", mstorage.Token)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
