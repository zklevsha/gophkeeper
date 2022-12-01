package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/zklevsha/gophkeeper/internal/helpers"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

// upassCreate  creates UserPassword entry and sends it to server via gRPC
func upassCreate(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {

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

	// converting to Pdata
	upass := structs.UPass{
		Name:     pname,
		Username: username,
		Password: password,
		Tags:     tags}
	pdata, err := helpers.ToPdata("upass", upass, mstorage.MasterKey)
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
func upassGet(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {
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
	pname := inputSelect("Upass name: ", pnames)
	resp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{Pname: pname})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}

	cleaned, err := helpers.FromPdata(resp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode upass: %s\n", err.Error())
	}
	up := cleaned.(structs.UPass)

	upass_pretty, err := json.MarshalIndent(up, "", " ")
	if err != nil {
		log.Printf("ERROR cant encode upass JSON : %s\n", err.Error())
	} else {
		log.Println(string(upass_pretty))
	}

}

func upassUpdate(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Getting current upass
	listResponse, err := gclient.Pdata.ListPdata(ctx, &pb.ListPdataRequest{Ptype: "upass"})
	if err != nil {
		log.Printf("ERROR: cant list pdata: %s\n", err.Error())
		return
	}
	if len(listResponse.PdataEtnry) == 0 {
		log.Println("you dont have any upass entries")
		return
	}
	var pnames []string
	for _, e := range listResponse.PdataEtnry {
		pnames = append(pnames, e.Name)
	}
	pname := inputSelect("Pname to update: ", pnames)
	getResp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{Pname: pname})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}
	cleaned, err := helpers.FromPdata(getResp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode upass: %s\n", err.Error())
		return
	}
	up := cleaned.(structs.UPass)

	// Parsing input
	nameNew := getInput(fmt.Sprintf("Name: [%s]", up.Name), any, false)
	if nameNew == "" {
		nameNew = up.Name
	}
	usernameNew := getInput(fmt.Sprintf("Username: [%s]", up.Username), any, false)
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
	tagsJson, err := json.Marshal(up.Tags)
	if err != nil {
		log.Printf("ERROR: cant parse old tags: %s\n", err.Error())
		return
	}
	tagsStr := getInput(fmt.Sprintf("new tags[%s]", tagsJson), isTags, false)
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
	upNew := structs.UPass{Name: nameNew,
		Username: usernameNew,
		Password: passwordNew,
		Tags:     tagsNew}
	pdataNew, err := helpers.ToPdata("upass", upNew, mstorage.MasterKey)
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

func upassDelete(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of available pnames
	listResponse, err := gclient.Pdata.ListPdata(ctx, &pb.ListPdataRequest{Ptype: "upass"})
	if err != nil {
		log.Printf("ERROR: cant list pdata: %s\n", err.Error())
		return
	}
	if len(listResponse.PdataEtnry) == 0 {
		log.Println("you dont have any upass entries")
		return
	}
	entries := make(map[string]int64)
	var pnames []string
	for _, e := range listResponse.PdataEtnry {
		pnames = append(pnames, e.Name)
		entries[e.Name] = e.ID
	}

	pname := inputSelect("Upass name: ", pnames)
	_, err = gclient.Pdata.DeletePdata(ctx, &pb.DeletePdataRequest{PdataID: entries[pname]})
	if err != nil {
		log.Printf("ERROR: cant delete pdata: %s", err.Error())
		return
	}
	log.Printf("upass %s (%d) was deleted", pname, entries[pname])
}
