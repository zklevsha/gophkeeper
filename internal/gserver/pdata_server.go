package gserver

import (
	"context"

	"github.com/zklevsha/gophkeeper/internal/db"
	"github.com/zklevsha/gophkeeper/internal/pb"
)

type pdataServer struct {
	pb.UnimplementedPrivateDataServer
	db  db.Connector
	key string
}

// AddPdata adds private data to database
func (s *authServer) AddPdata(ctx context.Context, in *pb.AddPdataRequest) (*pb.AddPdataResponse, error) {

}
