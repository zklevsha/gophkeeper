package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/zklevsha/gophkeeper/internal/pb"
)

// PricateString represents user`s private string
type Pstring struct {
	Name   string            `json:"name"`
	String string            `json:"string"`
	Tags   map[string]string `json:"tags,omitempty"`
}

// pstringCreate creates private string and sends it to server
func pstringCreate(mstorage *MemStorage, ctx context.Context, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// parsing input
	name := getInput("entry name:", notEmpty, false)
	string := getInput("string:", notEmpty, false)
	tags, err := getTags(getInput(`metainfo: {"key":"value",...}`, isTags, false))
	if err != nil {
		log.Printf("ERROR: cant parse tags: %s\n", err.Error())
		return
	}

	// converting to Pdata
	pstring := Pstring{Name: name, String: string, Tags: tags}
	pdata, err := toPdata("pstring", pstring, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cannot convert pstring to pdata: %s\n", err.Error())
		return
	}

	// sending pdata to server
	resp, err := gclient.Pdata.AddPdata(ctx, &pb.AddPdataRequest{Pdata: pdata})
	if err != nil {
		log.Printf("ERROR: cannot send pdata to server: %s\n", err.Error())
		return
	}
	log.Println(resp.Response)
}

// pstringGet retrives private string and sends it to server
func pstringGet(mstorage *MemStorage, ctx context.Context, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of existing pstring entries
	entries, err := listPnames(ctx, gclient, "pstring")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing pstring entries: %s", err.Error())
	}
	if len(entries) == 0 {
		log.Println("You dont have any pstring entries")
		return
	}
	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}

	// parsing input
	pname := inputSelect("Pstring name: ", pnames)
	pdataID := entries[pname]
	resp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{PdataID: pdataID})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}

	cleaned, err := fromPdata(resp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode pstring: %s\n", err.Error())
		return
	}
	pstring := cleaned.(Pstring)

	pstring_pretty, err := json.MarshalIndent(pstring, "", " ")
	if err != nil {
		log.Printf("ERROR cant encode pstring JSON : %s\n", err.Error())
	} else {
		log.Println(string(pstring_pretty))
	}

}

// pstringUpdate updates pstring entry
func pstringUpdate(mstorage *MemStorage, ctx context.Context, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Getting current pstring
	entries, err := listPnames(ctx, gclient, "pstring")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing pstring entries: %s", err.Error())
	}
	if len(entries) == 0 {
		log.Printf("You dont have any pstring entries")
		return
	}

	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}
	pname := inputSelect("Pstring to update: ", pnames)
	pdataID := entries[pname]
	getResp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{PdataID: pdataID})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}
	cleaned, err := fromPdata(getResp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode pstring: %s\n", err.Error())
		return
	}
	pstring := cleaned.(Pstring)

	// Parsing input
	nameNew := getInput(fmt.Sprintf("Name [%s]:", pstring.Name), any, false)
	if nameNew == "" {
		nameNew = pstring.Name
	}
	stringNew := getInput(fmt.Sprintf("String [%s]:", pstring.String), any, false)
	if stringNew == "" {
		stringNew = pstring.String
	}

	tagsJson, err := json.Marshal(pstring.Tags)
	if err != nil {
		log.Printf("ERROR: cant parse old tags: %s\n", err.Error())
		return
	}
	tagsStr := getInput(fmt.Sprintf("New tags [%s]", tagsJson), isTags, false)
	var tagsNew map[string]string
	if tagsStr != "" {
		tagsNew, err = getTags(tagsStr)
		if err != nil {
			log.Printf("ERROR: cant convert tags %s\n", err.Error())
			return
		}
	} else {
		tagsNew = pstring.Tags
	}

	// Convering to pdata
	pstringNew := Pstring{
		Name:   nameNew,
		String: stringNew,
		Tags:   tagsNew}
	pdataNew, err := toPdata("pstring", pstringNew, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cannot convert new pstring to pdata: %s\n", err.Error())
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
		log.Printf("ERROR: cant update pstring: %s\n", err.Error())
		return
	}
	log.Println(updateResp.Response)
}

// pstringDelete deletes pstring entry
func pstringDelete(mstorage *MemStorage, ctx context.Context, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of available pnames
	entries, err := listPnames(ctx, gclient, "pstring")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing pstring entries: %s", err.Error())
	}
	if len(entries) == 0 {
		log.Printf("You dont have any pstring entries")
		return
	}
	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}

	pname := inputSelect("pstring name: ", pnames)
	if !getYN(fmt.Sprintf("do you want delete %s?", pname)) {
		log.Println("Canceled")
		return
	}
	_, err = gclient.Pdata.DeletePdata(ctx, &pb.DeletePdataRequest{PdataID: entries[pname]})
	if err != nil {
		log.Printf("ERROR: cant delete pdata: %s", err.Error())
		return
	}
	log.Printf("pstring %s was deleted", pname)
}
