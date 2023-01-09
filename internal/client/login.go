package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/zklevsha/gophkeeper/internal/pb"
)

func login(ctx context.Context, gclient *Gclient,

	mstorage *MemStorage) error {
	email, err := getInput("email:", isEmail, false)
	if err != nil {
		return err
	}
	password, err := getInput("password:", notEmpty, true)
	if err != nil {
		return err
	}
	user := pb.User{Email: email, Password: password}
	ctxChild, cancel := context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	resp, err := gclient.Auth.GetToken(ctxChild, &pb.GetTokenRequest{User: &user})
	if err != nil {
		return fmt.Errorf("cant get token: %s", err.Error())
	}
	mstorage.SetToken(resp.Token)
	return nil
}

func register(ctx context.Context, gclient *Gclient) error {
	email, err := getInput("email:", isEmail, false)
	if err != nil {
		return err
	}
	password, err := getInput("password:", notEmpty, true)
	if err != nil {
		return err
	}
	paswordConfirm, err := getInput("password(confirm):", notEmpty, true)
	if err != nil {
		return err
	}
	if password != paswordConfirm {
		return errors.New("password mismatch")
	}
	user := pb.User{Email: email, Password: password}
	ctxChild, cancel := context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	_, err = gclient.Auth.Register(ctxChild, &pb.RegisterRequest{User: &user})
	if err != nil {
		return fmt.Errorf("cant register %s", err.Error())
	}
	return nil
}
