package client

import (
	"context"

	"github.com/zklevsha/gophkeeper/internal/structs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// list of RPC that not required authorization
var noAuth = map[string]bool{
	"/Auth/Register": true,
	"/Auth/GetToken": true,
}

// GetUnaryClientInterceptor returns a client interceptor
func GetUnaryClientInterceptor(mstorage *structs.MemStorage) grpc.UnaryClientInterceptor {
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
