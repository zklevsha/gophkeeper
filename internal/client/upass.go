package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

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
	pname, err := getInput("entry name:", notEmpty, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	username, err := getInput("username:", notEmpty, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	passwordOne, err := getInput("password (set empty for automatic generation):",
		any, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	if passwordOne == "" {
		password = GetRandomSrt(32)
	} else {
		passwordTwo, err := getInput("confirm password:", notEmpty, false)
		if err != nil {
			log.Printf("ERROR: parsing failed: %s", err.Error())
			return
		}
		if passwordOne != passwordTwo {
			log.Println("ERROR: password mismatch")
			return
		}
		password = passwordOne
	}
	tagsRaw, err  := getInput(`metainfo: {"key":"value",...}`, isTags, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	tags, err := getTags(tagsRaw)
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
	pdata, err := ToPdata("upass", upass, mstorage.MasterKey)
	if err != nil {
		log.Printf("canntot convert to Pdata: %s", err.Error())
	}

	// sending data to server
	ctxChild, cancel := context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	resp, err := gclient.Pdata.AddPdata(ctxChild, &pb.AddPdataRequest{Pdata: pdata})
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
	ctxChild, cancel := context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	entries, err := listPnames(ctxChild, gclient, "upass")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing upass entries: %s", err.Error())
		return
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
	pname, err := inputSelect("Upass name ", pnames)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	pdataID := entries[pname]

	// sending request to server
	ctxChild, cancel = context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	resp, err := gclient.Pdata.GetPdata(ctxChild, &pb.GetPdataRequest{PdataID: pdataID})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}

	cleaned, err := FromPdata(resp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode upass: %s\n", err.Error())
		return
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
	ctxChild, cancel := context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	entries, err := listPnames(ctxChild, gclient, "upass")
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
	pname, err := inputSelect("Pname to update: ", pnames)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	pdataID := entries[pname]
	ctxChild, cancel = context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	getResp, err := gclient.Pdata.GetPdata(ctxChild, &pb.GetPdataRequest{PdataID: pdataID})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}
	cleaned, err := FromPdata(getResp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode upass: %s\n", err.Error())
		return
	}
	up := cleaned.(UPass)

	// Parsing input
	nameNew, err := getInput(fmt.Sprintf("Name [%s]:", up.Name), any, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	if nameNew == "" {
		nameNew = up.Name
	}
	usernameNew, err := getInput(fmt.Sprintf("Username [%s]:", up.Username), any, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	if usernameNew == "" {
		usernameNew = up.Username
	}
	passwordNew, err := getInput(fmt.Sprintf("Password [%s]:", up.Password), any, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	if passwordNew == "" {
		passwordNew = up.Password
	} else {
		passwordNewConfirm, err := getInput("confirm password:", notEmpty, false)
		if err != nil {
			log.Printf("ERROR: parsing failed: %s", err.Error())
			return
		}
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
	tagsStr, err := getInput(fmt.Sprintf("New tags [%s]", tagsJSON), isTags, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
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
	pdataNew, err := ToPdata("upass", upNew, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cannot convert new upass to pdata: %s\n", err.Error())
		return
	}
	pdataNew.ID = getResp.Pdata.ID

	// Sending pdata to server
	ctxChild, cancel = context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	updateResp, err := gclient.Pdata.UpdatePdata(ctxChild, &pb.UpdatePdataRequest{Pdata: pdataNew})
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
	ctxChild, cancel := context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	entries, err := listPnames(ctxChild, gclient, "upass")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing upass entries: %s", err.Error())
		return
	}
	if len(entries) == 0 {
		log.Printf("You dont have any upass entries")
		return
	}
	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}

	pname, err := inputSelect("Upass name ", pnames)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	if !getYN(fmt.Sprintf("do you want delete %s?", pname)) {
		log.Println("Canceled")
		return
	}
	ctxChild, cancel = context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	_, err = gclient.Pdata.DeletePdata(ctxChild, &pb.DeletePdataRequest{PdataID: entries[pname]})
	if err != nil {
		log.Printf("ERROR: cant delete pdata: %s", err.Error())
		return
	}
	log.Printf("upass %s was deleted", pname)
}
