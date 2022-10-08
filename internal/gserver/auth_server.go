package gserver

import (
	"context"
	"errors"
	"fmt"

	"github.com/zklevsha/gophkeeper/internal/db"
	"github.com/zklevsha/gophkeeper/internal/jwt"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
	"golang.org/x/crypto/bcrypt"
)

type authServer struct {
	pb.UnimplementedAuthServer
	db  db.Connector
	key string
}

// Register register user it the system
func (s *authServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	// Decode
	if in.User == nil {
		response := pb.Response{Message: "", Error: "user is not set"}
		return &pb.RegisterResponse{Response: &response}, nil
	}

	encPass, err := bcrypt.GenerateFromPassword([]byte(in.User.Password), 14)
	if err != nil {
		e := fmt.Sprintf("failed to generate hash: %s", err.Error())
		response := pb.Response{Message: "", Error: e}
		return &pb.RegisterResponse{Response: &response}, nil
	}

	user := structs.User{Email: in.User.Email, Password: string(encPass)}

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

// GetToken authenticates user and generates JWT token
func (s *authServer) GetToken(ctx context.Context, in *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	// Decode
	if in.User == nil {
		response := pb.Response{Message: "", Error: "user is not set"}
		return &pb.GetTokenResponse{Response: &response}, nil
	}

	// authenticate
	user, err := s.db.GetUser(in.User.Email)
	if err != nil {
		var response pb.Response
		if errors.Is(err, structs.ErrUserAuth) {
			response = pb.Response{Message: "", Error: "authentication error"}
			return &pb.GetTokenResponse{Response: &response}, nil
		} else {
			response = pb.Response{Message: "",
				Error: fmt.Sprintf("db access errors: %s", err.Error())}
		}
		return &pb.GetTokenResponse{Response: &response}, nil
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(in.User.Password))
	if err != nil {
		response := pb.Response{Message: "", Error: "authentication error"}
		return &pb.GetTokenResponse{Response: &response}, nil
	}

	// generate JWT
	token, err := jwt.Generate(user.Id, s.key)
	if err != nil {
		e := fmt.Sprintf("cant generate token: %s", err.Error())
		response := pb.Response{Message: "", Error: e}
		return &pb.GetTokenResponse{Response: &response}, nil
	}
	response := pb.Response{Message: token, Error: ""}
	return &pb.GetTokenResponse{Response: &response}, nil
}
