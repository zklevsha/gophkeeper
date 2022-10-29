package client

import (
	"fmt"

	"github.com/zklevsha/gophkeeper/internal/enc"
	"github.com/zklevsha/gophkeeper/internal/helpers"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

func UPassCreate(mstorage *structs.MemStorage) error {
	var password string
	// parsing input
	// username
	username := getInput("username:", notEmpty, false)
	// password
	passwordOne := getInput("password (set empty for automatic generation):",
		any, false)
	if password == "" {
		password = helpers.GetRandomSrt(32)
	} else {
		passwordTwo := getInput("confirm password:", notEmpty, false)
		if passwordOne != passwordTwo {
			return fmt.Errorf("password mismatch")
		}
	}
	password = passwordOne
	// meta
	tags, err := getTags(getInput(`metainfo: {"key":"value",...}`, any, false))
	if err != nil {
		return fmt.Errorf("cant parse tags: %s", err.Error())
	}
	upass := structs.UPass{Username: username, Password: password, Tags: tags}

	// encrypting
	enc_upass, err := enc.EncryptAES([]byte(upass.Password), []byte(mstorage.MasterKey.Key))
	if err != nil {
		return fmt.Errorf("failed to encrypt input: %s", err.Error())
	}

}
