package gserver

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/zklevsha/gophkeeper/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)





func TestRegister(t *testing.T) {
	ctx := context.Background()
	server := setUp()
	defer tearDown(server)

	tt := []struct {
		name  string
		req *pb.RegisterRequest
	}{
		{name: "Registering new user",
			req: &pb.RegisterRequest{User: &pb.User{Email: "vasya@test.ru", Password: "secret"}},
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := server.auth.Register(ctx, tc.req)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}

}

func TestGetToken(t *testing.T) {
	server := setUp()
	defer tearDown(server)
	ctx := context.Background()
	testUser := pb.User{Email: "vasya@test.ru", Password: "secret"}
	// adding  test user
	_, err := server.auth.Register(ctx, &pb.RegisterRequest{User: &testUser})
	if err != nil {
		log.Fatalf("Cant register a test user: %s", err.Error())
	}

	tt := []struct {
		name  string
		user *pb.User
		errWant error
	}{
		{
			name: "Good user",
			user: &testUser ,
			errWant: nil,
		},
		{
			name: "Bad user",
			user: &pb.User{Email: "bad@bad.ru", Password: "bad"},
			errWant: status.Errorf(codes.Unauthenticated, "authentication failed"),
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			_, err := server.auth.GetToken(ctx, &pb.GetTokenRequest{User: tc.user})
			if !errors.Is(err, tc.errWant) {
				t.Errorf("error mismatch: have %s, want %s", err, tc.errWant)
			}
		})
	}

}