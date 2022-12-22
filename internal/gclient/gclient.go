package gclient

import (
	"google.golang.org/grpc"

	"github.com/zklevsha/gophkeeper/internal/pb"
)

// Gclient represents collection of various gRPC clients
type Gclient struct {
	Auth  pb.AuthClient
	Pdata pb.PrivateDataClient
}

func NewGclient(conn *grpc.ClientConn) Gclient {
	return Gclient{
		Auth:  pb.NewAuthClient(conn),
		Pdata: pb.NewPrivateDataClient(conn)}
}
