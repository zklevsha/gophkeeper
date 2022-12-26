package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/zklevsha/gophkeeper/internal/pb"
)

// Gclient represents collection of various gRPC clients
type Gclient struct {
	Auth  pb.AuthClient
	Pdata pb.PrivateDataClient
}

// NewGclient initializes new Gclient
func NewGclient(conn *grpc.ClientConn) Gclient {
	return Gclient{
		Auth:  pb.NewAuthClient(conn),
		Pdata: pb.NewPrivateDataClient(conn)}
}


// list of RPC that not required authorization
var noAuth = map[string]bool{
	"/Auth/Register": true,
	"/Auth/GetToken": true,
}

// GetUnaryClientInterceptor returns a client interceptor
func GetUnaryClientInterceptor(mstorage *MemStorage) grpc.UnaryClientInterceptor {
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
