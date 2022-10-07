package gserver

import (
	"github.com/zklevsha/gophkeeper/internal/db"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GetServer registers all gRPC services and returns grpc.Server instanse
func GetServer(db db.Connector) *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	pb.RegisterAuthServer(s, &authServer{db: db})
	return s

}
