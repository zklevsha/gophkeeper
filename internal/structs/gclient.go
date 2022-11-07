package structs

import (
	"github.com/zklevsha/gophkeeper/internal/pb"
	"google.golang.org/grpc"
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
