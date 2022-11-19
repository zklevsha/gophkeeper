package client

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/zklevsha/gophkeeper/internal/enc"
	"github.com/zklevsha/gophkeeper/internal/helpers"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

// upassCreate  creates UserPassword entry and sends it to server via gRPC
func upassCreate(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {

	err := upassCreateGetCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
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
			log.Println("ERROR: password mismatch")
			return
		}
		password = passwordOne
	}

	tags, err := getTags(getInput(`metainfo: {"key":"value",...}`, isTags, false))
	if err != nil {
		log.Printf("ERROR: cant parse tags: %s\n", err.Error())
		return
	}
	upass := structs.UPass{Username: username, Password: password, Tags: tags}

	// encoding
	upass_encoded, err := json.Marshal(upass)
	if err != nil {
		log.Printf("ERROR: cannot encode upass to JSON: %s\n", err.Error())
		return
	}

	// encrypting
	upass_encrypted, err := enc.EncryptAES(upass_encoded, []byte(mstorage.MasterKey.Key))
	if err != nil {
		log.Printf("ERROR: failed to encrypt upass: %s\n", err.Error())
		return
	}
	// sending to pdata to server
	pdata := pb.Pdata{
		Pname:   pname,
		Ptype:   "upass",
		Pdata:   upass_encrypted,
		KeyHash: mstorage.MasterKey.KeyHash[:]}
	resp, err := gclient.Pdata.AddPdata(ctx, &pb.AddPdataRequest{Pdata: &pdata})
	if err != nil {
		log.Printf("ERROR: cant send message to server: %s\n", err.Error())
		return
	}
	log.Println(resp.Response)

}

// upassCheck checks that memstorage has all needed infomation
// for upass creation/retrival (token, masterkey...)
func upassCreateGetCheck(mstorage *structs.MemStorage) error {
	if mstorage.MasterKey.Key == "" {
		return errors.New("master-key does not exists add it via key-generate/key-load commands")
	}
	if mstorage.Token == "" {
		return errors.New("login required (login)")
	}
	return nil
}

// upassGet retrives Upass from gRPC server
func upassGet(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {
	err := upassCreateGetCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	pname := getInput("upass name:", notEmpty, false)
	resp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{Pname: pname})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}

	// check if we using correct master key
	if string(resp.Pdata.KeyHash) != string(mstorage.MasterKey.KeyHash) {
		log.Println("ERROR: key hash mismatch")
		log.Printf("hash of the key that used to encrypt pdata: %v\n",
			resp.Pdata.KeyHash)
		log.Printf("masterkey hash: %v\n", mstorage.MasterKey.KeyHash)
		return
	}

	// decrypt
	upass_decrypted, err := enc.DecryptAES(resp.Pdata.Pdata,
		[]byte(mstorage.MasterKey.Key))
	if err != nil {
		log.Printf("ERROR cant decrypt data: %s\n", err.Error())
		return
	}

	// decode
	var up structs.UPass
	err = json.Unmarshal(upass_decrypted, &up)
	if err != nil {
		log.Printf("ERROR cant decode upass JSON to struct: %s\n", err.Error())
		return
	}
	upass_pretty, err := json.MarshalIndent(up, "", " ")
	if err != nil {
		log.Printf("ERROR cant encode upass JSON : %s\n", err.Error())
	} else {
		log.Println(string(upass_pretty))
	}

}
