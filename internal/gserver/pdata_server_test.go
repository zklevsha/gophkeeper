package gserver

import (
	"context"
	"testing"

	"github.com/zklevsha/gophkeeper/internal/client"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"google.golang.org/grpc/metadata"
)

// adding user to server and adding JWT to context
func getTokenCtx(ctx context.Context, user *pb.User, s testServer) (context.Context, error){
	_, err := s.auth.Register(ctx, &pb.RegisterRequest{User: user})
	if err != nil {
		return nil, err
	}
	// geting token
	getTokenResp, err := s.auth.GetToken(ctx, &pb.GetTokenRequest{User: user})
	if err != nil {
		return nil, err
	}
	// adding token to ctx
	md := metadata.MD{}
	md.Append("authorization", getTokenResp.Token)
	md.Append("test", "test")
	ctx = metadata.NewIncomingContext(ctx, md)

	return ctx, nil
}

func TestGetUserId(t *testing.T) {
	server := setUp()
	defer tearDown(server)

	var want int64 = 1
	testUser := pb.User{Email: "vasya@test.ru", Password: "secret"}
	ctx := context.Background()

	ctx, err := getTokenCtx(ctx, &testUser, server)
	if err != nil {
		t.Fatalf("cant get context with ctx: %s",err.Error())
	}

	have, err := server.pdata.getUserID(ctx)
	if err != nil {
		t.Fatalf("getUserID have returned error: %s", err.Error())
	}

	if have != want  {
		t.Errorf("return mismatch: have %d, want  %d", have, want)
	}

}

func TestPdata(t *testing.T) {
	server := setUp()
	defer tearDown(server)
	testUser := pb.User{Email: "vasya@test.ru", Password: "secret"}
	ctx := context.Background()

	ctx, err := getTokenCtx(ctx, &testUser, server)
	if err != nil {
		t.Fatalf("cant get context with ctx: %s",err.Error())
	}

	key := client.MasterKey{Key: client.GetRandomSrt(32)}
	testPdata, err := client.ToPdata("pstring", client.Pstring{Name: "test", String: "secret"}, key)
	if err != nil {
		t.Fatalf("cant convert test data to Pdata: %s", err.Error())
	}

	_, err = server.pdata.AddPdata(ctx, &pb.AddPdataRequest{Pdata: testPdata})
	if err != nil {
		t.Errorf("AddPdata have returned error: %s", err.Error())
	}


	testPdata.ID = 1
	_, err = server.pdata.UpdatePdata(ctx, &pb.UpdatePdataRequest{Pdata: testPdata})
	if err != nil {
		t.Errorf("UpdatePdata have returned error: %s", err.Error())
	}

	_, err = server.pdata.DeletePdata(ctx, &pb.DeletePdataRequest{PdataID: 1})
	if err != nil {
		t.Fatalf("cant delete test data %s", err.Error())
	}
}