// authServer implements Register and Authentication logic

package gserver

import (
	"context"
	"fmt"

	"github.com/zklevsha/gophkeeper/internal/client"
	"github.com/zklevsha/gophkeeper/internal/db"
	"github.com/zklevsha/gophkeeper/internal/jmanager"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return nil, status.Errorf(codes.InvalidArgument, "user is not set")
	}

	encPass, err := bcrypt.GenerateFromPassword([]byte(in.User.Password), 14)
	if err != nil {
		e := fmt.Sprintf("failed to generate hash: %s", err.Error())
		return nil, status.Errorf(codes.Internal, e)
	}

	user := client.User{Email: in.User.Email, Password: string(encPass)}

	id, err := s.db.Register(ctx, user)
	if err != nil {
		e := fmt.Sprintf("failed to register user: %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}
	r := fmt.Sprintf("user %s (userid: %d) was created", user.Email, id)
	return &pb.RegisterResponse{Response: r}, nil
}

// GetToken authenticates user and generates JWT token
func (s *authServer) GetToken(ctx context.Context, in *pb.GetTokenRequest) (*pb.GetTokenResponse, error) {
	// Decode
	if in.User == nil {
		return nil, status.Errorf(codes.InvalidArgument, "user is not set")
	}

	// authenticate
	user, err := s.db.GetUser(ctx, in.User.Email)
	if err != nil {
		return nil, status.Errorf(getCode(err), err.Error())
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(in.User.Password))
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication error")
	}

	// generate JWT
	token, err := jmanager.Generate(user.ID, s.key)
	if err != nil {
		e := fmt.Sprintf("cant generate token: %s", err.Error())
		return nil, status.Errorf(codes.Internal, e)
	}
	return &pb.GetTokenResponse{Token: token.Token}, nil
}
