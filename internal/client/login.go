package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/zklevsha/gophkeeper/internal/pb"
)

func login(ctx context.Context, gclient *Gclient,
	mstorage *MemStorage) error {
	email := getInput("email:", isEmail, false)
	password := getInput("password:", notEmpty, true)

	user := pb.User{Email: email, Password: password}
	resp, err := gclient.Auth.GetToken(ctx, &pb.GetTokenRequest{User: &user})
	if err != nil {
		return fmt.Errorf("cant get token: %s", err.Error())
	}
	mstorage.SetToken(resp.Token)
	return nil
}

func register(ctx context.Context, gclient *Gclient) error {
	email := getInput("email:", isEmail, false)
	password := getInput("password:", notEmpty, true)
	paswordConfirm := getInput("password(confirm):", notEmpty, true)
	if password != paswordConfirm {
		return errors.New("password mismatch")
	}
	user := pb.User{Email: email, Password: password}
	_, err := gclient.Auth.Register(ctx, &pb.RegisterRequest{User: &user})
	if err != nil {
		return fmt.Errorf("cant register %s", err.Error())
	}
	return nil
}
