package gserver

import (
	"context"
	"fmt"

	"github.com/zklevsha/gophkeeper/internal/db"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

type authServer struct {
	pb.UnimplementedAuthServer
	db db.Connector
}

// Register register user it the system
func (s *authServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Decode
	if in.User == nil {
		response := pb.Response{Message: "", Error: "user is not set"}
		return &pb.RegisterResponse{Response: &response}, nil
	}
	user := structs.User{Email: in.User.Email, Password: in.User.Password}

	id, err := s.db.Register(user)
	if err != nil {
		e := fmt.Sprintf("failed to register user: %s", err.Error())
		response := pb.Response{Message: "", Error: e}
		return &pb.RegisterResponse{Response: &response}, nil
	}
	m := fmt.Sprintf("user %s (userid: %d) was created", user.Email, id)
	response := pb.Response{Message: m, Error: ""}
	return &pb.RegisterResponse{Response: &response}, nil
}
