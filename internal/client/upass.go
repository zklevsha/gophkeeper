package client

import (
	"context"
	"errors"
	"fmt"

	"github.com/zklevsha/gophkeeper/internal/enc"
	"github.com/zklevsha/gophkeeper/internal/helpers"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

// upassCreate  creates UserPassword entry and sends it to server via gRPC
func upassCreate(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {

	err := upassCreateCheck(mstorage)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// parsing input
	var password string
	pname := getInput("entry name:", notEmpty, false)
	username := getInput("username:", notEmpty, false)
	passwordOne := getInput("password (set empty for automatic generation):",
		any, false)
	if passwordOne == "" {
		password = helpers.GetRandomSrt(32)
	} else {
		passwordTwo := getInput("confirm password:", notEmpty, false)
		if passwordOne != passwordTwo {
			fmt.Println("ERROR: password mismatch")
			return
		}
		password = passwordOne
	}

	tags, err := getTags(getInput(`metainfo: {"key":"value",...}`, any, false))
	if err != nil {
		fmt.Printf("ERROR: cant parse tags: %s\n", err.Error())
		return
	}
	upass := structs.UPass{Username: username, Password: password, Tags: tags}
	fmt.Println("1")
	// encrypting
	enc_upass, err := enc.EncryptAES([]byte(upass.Password), []byte(mstorage.MasterKey.Key))
	if err != nil {
		fmt.Printf("ERROR: failed to encrypt input: %s\n", err.Error())
		return
	}
	// sending to pdata to server
	pdata := pb.Pdata{
		Pname:   pname,
		Ptype:   "upass",
		Pdata:   enc_upass,
		KeyHash: mstorage.MasterKey.KeyHash[:]}
	resp, err := gclient.Pdata.AddPdata(ctx, &pb.AddPdataRequest{Pdata: &pdata})
	if err != nil {
		fmt.Printf("ERROR: cant send message to server: %s\n", err.Error())
		return
	}
	fmt.Println(resp.Response.Message)

}

// upassCheck checks that memstorage has all needed infomation
// for upass creation token, masterkey etc)

func upassCreateCheck(mstorage *structs.MemStorage) error {
	if mstorage.MasterKey.Key == "" {
		return errors.New("master-key does not exists add it via key-generate/key-load commands")
	}
	if mstorage.Token == "" {
		return errors.New("login required (login)")
	}
	return nil
}
