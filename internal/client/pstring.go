package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/zklevsha/gophkeeper/internal/pb"
)

// Pstring represents user`s private string
type Pstring struct {
	Name   string            `json:"name"`
	String string            `json:"string"`
	Tags   map[string]string `json:"tags,omitempty"`
}

// pstringCreate creates private string and sends it to server
func pstringCreate(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// parsing input
	name, err := getInput("entry name:", notEmpty, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	string, err := getInput("string:", notEmpty, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	tagsRaw, err := getInput(`metainfo: {"key":"value",...}`, isTags, false)
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
	pstring := Pstring{Name: name, String: string, Tags: tags}
	pdata, err := ToPdata("pstring", pstring, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cannot convert pstring to pdata: %s\n", err.Error())
		return
	}

	// sending pdata to server
	ctxChild, cancel := context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	resp, err := gclient.Pdata.AddPdata(ctxChild, &pb.AddPdataRequest{Pdata: pdata})
	if err != nil {
		log.Printf("ERROR: cannot send pdata to server: %s\n", err.Error())
		return
	}
	log.Println(resp.Response)
}

// pstringGet retrives private string and sends it to server
func pstringGet(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of existing pstring entries
	entries, err := listPnames(ctx, gclient, "pstring")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing pstring entries: %s", err.Error())
		return
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
	pname, err := inputSelect("Pstring name: ", pnames)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	pdataID := entries[pname]
	resp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{PdataID: pdataID})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}

	cleaned, err := FromPdata(resp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode pstring: %s\n", err.Error())
		return
	}
	pstring := cleaned.(Pstring)

	pstringPretty, err := json.MarshalIndent(pstring, "", " ")
	if err != nil {
		log.Printf("ERROR cant encode pstring JSON : %s\n", err.Error())
	} else {
		log.Println(string(pstringPretty))
	}

}

// pstringUpdate updates pstring entry
func pstringUpdate(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// Getting current pstring
	ctxChild, cancel := context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	entries, err := listPnames(ctxChild, gclient, "pstring")
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
	pname, err := inputSelect("Pstring to update: ", pnames)
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
		log.Printf("ERROR: cant decode pstring: %s\n", err.Error())
		return
	}
	pstring := cleaned.(Pstring)

	// Parsing input
	nameNew, err := getInput(fmt.Sprintf("Name [%s]:", pstring.Name), any, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	if nameNew == "" {
		nameNew = pstring.Name
	}
	stringNew, err := getInput(fmt.Sprintf("String [%s]:", pstring.String), any, false)
	if err != nil {
		log.Printf("ERROR: parsing failed: %s", err.Error())
		return
	}
	if stringNew == "" {
		stringNew = pstring.String
	}
	tagsJSON, err := json.Marshal(pstring.Tags)
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
		tagsNew = pstring.Tags
	}

	// Convering to pdata
	pstringNew := Pstring{
		Name:   nameNew,
		String: stringNew,
		Tags:   tagsNew}
	pdataNew, err := ToPdata("pstring", pstringNew, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cannot convert new pstring to pdata: %s\n", err.Error())
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
		log.Printf("ERROR: cant update pstring: %s\n", err.Error())
		return
	}
	log.Println(updateResp.Response)
}

// pstringDelete deletes pstring entry
func pstringDelete(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of available pnames
	ctxChild, cancel := context.WithTimeout(ctx, time.Duration(reqTimeout))
	defer cancel()
	entries, err := listPnames(ctxChild, gclient, "pstring")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing pstring entries: %s", err.Error())
		return
	}
	if len(entries) == 0 {
		log.Printf("You dont have any pstring entries")
		return
	}
	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}

	pname, err := inputSelect("pstring name: ", pnames)
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
	log.Printf("pstring %s was deleted", pname)
}
