package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/zklevsha/gophkeeper/internal/pb"
)

// UPass represents user`s Username/Password pair
type UPass struct {
	Name     string            `json:"name"`
	Username string            `json:"username"`
	Password string            `json:"password"`
	Tags     map[string]string `json:"tags,omitempty"`
}


// upassCreate  creates UserPassword entry and sends it to server via gRPC
func upassCreate(ctx context.Context,  mstorage *MemStorage, gclient *Gclient) {

	err := reqCheck(mstorage)
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
		password = getRandomSrt(32)
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

	// converting to Pdata
	upass := UPass{
		Name:     pname,
		Username: username,
		Password: password,
		Tags:     tags}
	pdata, err := toPdata("upass", upass, mstorage.MasterKey)
	if err != nil {
		log.Printf("canntot convert to Pdata: %s", err.Error())
	}

	// sending data to server
	resp, err := gclient.Pdata.AddPdata(ctx, &pb.AddPdataRequest{Pdata: pdata})
	if err != nil {
		log.Printf("ERROR: cant send message to server: %s\n", err.Error())
		return
	}
	log.Println(resp.Response)

}

// upassGet retrives Upass from gRPC server
func upassGet(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of existing upass entries
	entries, err := listPnames(ctx, gclient, "upass")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing upass entries: %s", err.Error())
	}
	if len(entries) == 0 {
		log.Println("You dont have any upass entries")
		return
	}
	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}

	// parsing input
	pname := inputSelect("Upass name ", pnames)
	pdataID := entries[pname]
	resp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{PdataID: pdataID})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}

	cleaned, err := fromPdata(resp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode upass: %s\n", err.Error())
	}
	up := cleaned.(UPass)

	upassPretty, err := json.MarshalIndent(up, "", " ")
	if err != nil {
		log.Printf("ERROR cant encode upass JSON: %s\n", err.Error())
	} else {
		log.Println(string(upassPretty))
	}

}

func upassUpdate(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Getting current upass
	entries, err := listPnames(ctx, gclient, "upass")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing upass entries: %s", err.Error())
	}
	if len(entries) == 0 {
		log.Printf("You dont have any upass entries")
		return
	}

	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}
	pname := inputSelect("Pname to update: ", pnames)
	pdataID := entries[pname]
	getResp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{PdataID: pdataID})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}
	cleaned, err := fromPdata(getResp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode upass: %s\n", err.Error())
		return
	}
	up := cleaned.(UPass)

	// Parsing input
	nameNew := getInput(fmt.Sprintf("Name [%s]:", up.Name), any, false)
	if nameNew == "" {
		nameNew = up.Name
	}
	usernameNew := getInput(fmt.Sprintf("Username [%s]:", up.Username), any, false)
	if usernameNew == "" {
		usernameNew = up.Username
	}
	passwordNew := getInput(fmt.Sprintf("Password [%s]:", up.Password), any, false)
	if passwordNew == "" {
		passwordNew = up.Password
	} else {
		passwordNewConfirm := getInput("confirm password:", notEmpty, false)
		if passwordNew != passwordNewConfirm {
			log.Println("ERROR: password mismatch")
			return
		}
	}
	tagsJSON, err := json.Marshal(up.Tags)
	if err != nil {
		log.Printf("ERROR: cant parse old tags: %s\n", err.Error())
		return
	}
	tagsStr := getInput(fmt.Sprintf("New tags [%s]", tagsJSON), isTags, false)
	var tagsNew map[string]string
	if tagsStr != "" {
		tagsNew, err = getTags(tagsStr)
		if err != nil {
			log.Printf("ERROR: cant convert tags %s\n", err.Error())
			return
		}
	} else {
		tagsNew = up.Tags
	}

	// Convering to pdata
	upNew := UPass{Name: nameNew,
		Username: usernameNew,
		Password: passwordNew,
		Tags:     tagsNew}
	pdataNew, err := toPdata("upass", upNew, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cannot convert new upass to pdata: %s\n", err.Error())
		return
	}
	pdataNew.ID = getResp.Pdata.ID

	// Sending pdata to server
	updateResp, err := gclient.Pdata.UpdatePdata(ctx, &pb.UpdatePdataRequest{Pdata: pdataNew})
	if err != nil {
		log.Printf("ERROR: cant send message to server: %s\n", err.Error())
		return
	}

	if err != nil {
		log.Printf("ERROR: cant update upass: %s\n", err.Error())
		return
	}
	log.Println(updateResp.Response)
}

func upassDelete(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of available pnames
	entries, err := listPnames(ctx, gclient, "upass")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing upass entries: %s", err.Error())
	}
	if len(entries) == 0 {
		log.Printf("You dont have any upass entries")
		return
	}
	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}

	pname := inputSelect("Upass name ", pnames)
	if !getYN(fmt.Sprintf("do you want delete %s?", pname)) {
		log.Println("Canceled")
		return
	}
	_, err = gclient.Pdata.DeletePdata(ctx, &pb.DeletePdataRequest{PdataID: entries[pname]})
	if err != nil {
		log.Printf("ERROR: cant delete pdata: %s", err.Error())
		return
	}
	log.Printf("upass %s was deleted", pname)
}
