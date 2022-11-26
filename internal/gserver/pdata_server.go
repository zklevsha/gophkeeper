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
// (userid is set by server interceptor)
func (s *pdataServer) getUserId(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.Internal, "failed to retrive metadata from ctx")
	}
	token_raw, ok := md["authorization"]
	if !ok {
		return 0, structs.ErrNoToken
	}
	token_parsed, err := jmanager.Validate(token_raw[0], s.key)
	if err != nil {
		return 0, structs.ErrInvalidToken
	}

	return token_parsed.Claims.UserID, nil

}

// AddPdata adds private data to database
func (s *pdataServer) AddPdata(ctx context.Context, in *pb.AddPdataRequest) (*pb.AddPdataResponse, error) {
	if in.Pdata == nil {
		return nil, status.Errorf(codes.InvalidArgument, "pdata is nil")
	}

	pdata := structs.Pdata{
		Name:        in.Pdata.Pname,
		Type:        in.Pdata.Ptype,
		KeyHash:     base64.StdEncoding.EncodeToString(in.Pdata.KeyHash),
		PrivateData: base64.StdEncoding.EncodeToString(in.Pdata.Pdata)}

	userID, err := s.getUserId(ctx)
	if err != nil {
		e := fmt.Sprintf("failed to get userid: %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}
	_, err = s.db.PrivateAdd(userID, pdata)
	if err != nil {
		e := fmt.Sprintf("failed to add pdata to database: %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}
	r := fmt.Sprintf("pdata %s was added sucsessfully", pdata.Name)
	return &pb.AddPdataResponse{Response: r}, nil

}

func (s *pdataServer) GetPdata(ctx context.Context, in *pb.GetPdataRequest) (*pb.GetPdataResponse, error) {
	if in.Pname == "" {
		return nil, status.Errorf(codes.InvalidArgument, "pname is not set")
	}
	userID, err := s.getUserId(ctx)
	if err != nil {
		e := fmt.Sprintf("failed to get userid: %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}

	pdataID, err := s.db.GetPdataID(userID, in.Pname)
	if err != nil {
		e := fmt.Sprintf("failed to get pdataid: %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}

	pdata, err := s.db.PrivateGet(userID, pdataID)
	if err != nil {
		e := fmt.Sprintf("failed to get pdata: %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}

	keyHash, err := base64.StdEncoding.DecodeString(pdata.KeyHash)
	if err != nil {
		e := fmt.Sprintf("cant decode khash_base64 to byte array: %s", err.Error())
		return nil, status.Error(codes.Internal, e)
	}
	privateData, err := base64.StdEncoding.DecodeString(pdata.PrivateData)
	if err != nil {
		e := fmt.Sprintf("cant decode data_base64 to byte array: %s", err.Error())
		return nil, status.Error(codes.Internal, e)
	}
	pbPdata := pb.Pdata{
		ID:      pdataID,
		Pname:   pdata.Name,
		Ptype:   pdata.Type,
		KeyHash: keyHash,
		Pdata:   privateData}

	return &pb.GetPdataResponse{Pdata: &pbPdata}, nil
}

func (s *pdataServer) UpdatePdata(ctx context.Context, in *pb.UpdatePdataRequest) (*pb.UpdatePdataResponse, error) {
	if in.Pdata == nil {
		return nil, status.Errorf(codes.InvalidArgument, "pname is not set")
	}
	userID, err := s.getUserId(ctx)
	if err != nil {
		e := fmt.Sprintf("failed to get userid: %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}
	pdata := structs.Pdata{
		ID:          in.Pdata.ID,
		Name:        in.Pdata.Pname,
		Type:        in.Pdata.Ptype,
		KeyHash:     base64.StdEncoding.EncodeToString(in.Pdata.KeyHash),
		PrivateData: base64.StdEncoding.EncodeToString(in.Pdata.Pdata)}

	err = s.db.PrivateUpdate(userID, pdata)
	if err != nil {
		e := fmt.Sprintf("failed to update pdata: %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}
	r := fmt.Sprintf("pdata %s was updated sucsessfully", pdata.Name)
	return &pb.UpdatePdataResponse{Response: r}, nil
}

func (s *pdataServer) ListPdata(ctx context.Context, in *pb.ListPdataRequest) (*pb.ListPdataResponse, error) {
	userID, err := s.getUserId(ctx)
	if err != nil {
		e := fmt.Sprintf("failed to get userid: %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}

	pdata, err := s.db.PrivateList(userID, in.Ptype)
	if err != nil {
		e := fmt.Sprintf("failed to list pdata: %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}

	var pdataPointers []*pb.PdataEntry
	for _, p := range pdata {
		pdataPointers = append(pdataPointers,
			&pb.PdataEntry{Name: p.Name, ID: p.ID})
	}

	return &pb.ListPdataResponse{PdataEtnry: pdataPointers}, nil

}
