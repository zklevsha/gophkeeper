package client

import (
	"context"
	"fmt"
	"log"

	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

func login(ctx context.Context, gclient *structs.Gclient,
	mstorage *structs.MemStorage) {
	email := getInput("email:", isEmail, false)
	password := getInput("password:", notEmpty, true)

	user := pb.User{Email: email, Password: password}
	resp, err := gclient.Auth.GetToken(ctx, &pb.GetTokenRequest{User: &user})
	if err != nil {
		log.Printf("ERROR cant get token: %s\n", err.Error())
		return
	}
	mstorage.SetToken(resp.Token)
	fmt.Println("login successful")

}

func register(ctx context.Context, gclient *structs.Gclient) {
	email := getInput("email:", isEmail, false)
	password := getInput("password:", notEmpty, true)
	paswordConfirm := getInput("password(confirm):", notEmpty, true)
	if password != paswordConfirm {
		fmt.Println("ERROR password mismatch")
		return
	}
	user := pb.User{Email: email, Password: password}
	resp, err := gclient.Auth.Register(ctx, &pb.RegisterRequest{User: &user})
	if err != nil {
		log.Printf("ERROR cant register %s\n", err.Error())
		return
	}
	log.Printf("%s\n", resp.Response)
}
