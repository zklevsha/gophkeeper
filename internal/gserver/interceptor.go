package gserver

import (
	"context"
	"fmt"
	"log"

	"github.com/zklevsha/gophkeeper/internal/jmanager"
	"github.com/zklevsha/gophkeeper/internal/structs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// list of RPC that not require authorization
var noAuth = map[string]bool{
	"/Auth/Register": true,
	"/Auth/GetToken": true,
}

// GetUnaryServerInterceptor returns server unary Interceptor to authenticate and authorize unary RPC
func GetUnaryServerInterceptor(jwtKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log.Println("--> unary interceptor: ", info.FullMethod)

		if noAuth[info.FullMethod] {
			log.Println("---> unary interceptor: adding parsint access token")
			token, err := parseJTW(ctx, jwtKey)
			if err != nil {
				return nil, err
			}
			ctx = metadata.AppendToOutgoingContext(ctx, "userid", fmt.Sprintf("%d", token.Claims.UserID))
		}

		return handler(ctx, req)
	}
}

// parseJWT  checks JWT and parse it
func parseJTW(ctx context.Context, jwtKey string) (structs.Jtoken, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return structs.Jtoken{}, status.Errorf(codes.InvalidArgument, "Retrieving metadata is failed")
	}

	authHeader, ok := md["authorization"]
	if !ok {
		return structs.Jtoken{}, status.Errorf(codes.Unauthenticated, "Authorization token is not supplied")
	}

	token, err := jmanager.Validate(authHeader[0], jwtKey)

	if err != nil {
		return structs.Jtoken{}, status.Errorf(codes.Unauthenticated, err.Error())
	}
	return token, nil
}
