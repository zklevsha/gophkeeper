package gserver

import (
	"context"
	"testing"

	"github.com/zklevsha/gophkeeper/internal/pb"
	"google.golang.org/grpc/metadata"
)

func TestGetUserId(t *testing.T) {
	server := setUp()
	defer tearDown(server)

	var want int64 = 1
	testUser := pb.User{Email: "vasya@test.ru", Password: "secret"}
	ctx := context.Background()
	// adding  test user
	_, err := server.auth.Register(ctx, &pb.RegisterRequest{User: &testUser})
	if err != nil {
		t.Fatal(err.Error())
	}
	// geting token
	resp, err := server.auth.GetToken(ctx, &pb.GetTokenRequest{User: &testUser})
	if err != nil {
		t.Fatal(err.Error())
	}

	md := metadata.MD{}
	md.Append("authorization", resp.Token)
	md.Append("test", "test")
	ctx = metadata.NewIncomingContext(ctx, md)
	have, err := server.pdata.getUserID(ctx)
	if err != nil {
		t.Fatalf("getUserID have returned error: %s", err.Error())
	}

	if have != want  {
		t.Errorf("return mismatch: have %d, want  %d", have, want)
	}

}