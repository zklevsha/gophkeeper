package gserver

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"

	"github.com/zklevsha/gophkeeper/internal/db"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const dsnDefault = "postgres://gophkeeper:gophkeeper@localhost:5532/gophkeeper_test?sslmode=disable"

var ctx = context.Background()

func setUp() authServer {
	// added so github action will be able to connect to test database
	var dsn = os.Getenv("GK_DB_TEST_DSN")
	if (dsn == "" ){
		dsn = dsnDefault
	}
	// connecting to DB
	c := db.Connector{Ctx: ctx, DSN: dsn}
	err := c.Init()
	if err != nil {
		log.Fatalf("Failed to init Connector: %s", err.Error())
	}

	// running migrations
	err = db.RunMigrations(dsn, "up")
	if err != nil {
		log.Fatalf("cannot run up migrations: %s", err.Error())
	}

	return authServer{db: c, key: "secret"}
}

func tearDown(s authServer) {
	err := db.RunMigrations(s.db.DSN, "down")
	if err != nil {
		log.Fatalf("cannot run down migrations: %s", err.Error())
	}
	s.db.Close()
}





func TestRegister(t *testing.T) {
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
			_, err := server.Register(ctx, tc.req)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}

}

func TestGetToken(t *testing.T) {
	server := setUp()
	defer tearDown(server)
	testUser := pb.User{Email: "vasya@test.ru", Password: "secret"}
	// adding  test user
	_, err := server.Register(ctx, &pb.RegisterRequest{User: &testUser})
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
			_, err := server.GetToken(ctx, &pb.GetTokenRequest{User: tc.user})
			if !errors.Is(err, tc.errWant) {
				t.Errorf("error mismatch: have %s, want %s", err, tc.errWant)
			}
		})
	}

}