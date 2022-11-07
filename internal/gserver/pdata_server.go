package gserver

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/zklevsha/gophkeeper/internal/db"
	"github.com/zklevsha/gophkeeper/internal/jmanager"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type pdataServer struct {
	pb.UnimplementedPrivateDataServer
	db  db.Connector
	key string
}

// getUserId retrives userid from context
// userid supposed to be set by server interceptor
func (s *pdataServer) getUserId(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.Internal, "failed to retrive metadata from ctx")
	}
	fmt.Printf("%v", md)
	token_raw, ok := md["authorization"]
	if !ok {
		e := "failed to retrive token from context: no 'authorization' key in context metadata"
		return 0, status.Errorf(codes.PermissionDenied, e)
	}
	token_parsed, err := jmanager.Validate(token_raw[0], s.key)
	if err != nil {
		e := fmt.Sprintf("token validation failed: %s", err.Error())
		return 0, status.Error(codes.PermissionDenied, e)
	}

	return token_parsed.Claims.UserID, nil

}

// AddPdata adds private data to database
func (s *pdataServer) AddPdata(ctx context.Context, in *pb.AddPdataRequest) (*pb.AddPdataResponse, error) {
	if in.Pdata == nil {
		return nil, status.Errorf(codes.NotFound, "pdata is nil")
	}

	pdata := structs.Pdata{
		Name:        in.Pdata.Pname,
		Type:        in.Pdata.Ptype,
		KeyHash:     base64.StdEncoding.EncodeToString(in.Pdata.KeyHash),
		PrivateData: base64.StdEncoding.EncodeToString(in.Pdata.Pdata)}

	userID, err := s.getUserId(ctx)
	if err != nil {
		e := fmt.Sprintf("failed to get userid: %s", err.Error())
		return nil, status.Errorf(codes.Internal, e)
	}
	err = s.db.PrivateAdd(userID, pdata)
	if err != nil {
		e := fmt.Sprintf("failed to add pdata to database: %s", err.Error())
		return nil, status.Errorf(codes.Internal, e)
	}
	response := pb.Response{
		Message: fmt.Sprintf("pdata %s was added sucsessfully", pdata.Name), Error: ""}
	return &pb.AddPdataResponse{Response: &response}, nil

}
