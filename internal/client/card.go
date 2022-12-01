package client

import (
	"context"
	"encoding/json"
	"log"

	"github.com/zklevsha/gophkeeper/internal/helpers"
	"github.com/zklevsha/gophkeeper/internal/pb"
	"github.com/zklevsha/gophkeeper/internal/structs"
)

// cardCreate creates Credit card entry and sends it to server via gRPC
func cardCreate(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {

	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// parsing input
	name := getInput("name:", notEmpty, false)
	number := getInput("number(XXXX XXXX XXXX XXXX):", isCardNumber, false)
	holder := getInput("holder(JOHN DOE):", isCardHolder, false)
	expire := getInput("expire(MM/YY):", isCardExire, false)
	cvc := getInput("CVV/CVC:", isCardCVC, false)
	tags, err := getTags(getInput(`metainfo: {"key":"value",...}`, isTags, false))
	if err != nil {
		log.Printf("ERROR: cant parse tags: %s\n", err.Error())
		return
	}

	// converting to Pdata
	card := structs.Card{
		Name: name, Number: number,
		Holder: holder, Expire: expire,
		CVC: cvc, Tags: tags}
	pdata, err := helpers.ToPdata("card", card, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant convert to pdata %s", err.Error())
	}

	// sending data to server
	resp, err := gclient.Pdata.AddPdata(ctx, &pb.AddPdataRequest{Pdata: pdata})
	if err != nil {
		log.Printf("ERROR: cant send message to server: %s\n", err.Error())
		return
	}
	log.Println(resp.Response)
}

// cardCreate retreives Credit card entry from the server
func cardGet(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {

	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of existing card entries
	entries, err := listPnames(ctx, gclient, "card")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing upass entries: %s", err.Error())
	}
	if len(entries) == 0 {
		log.Printf("You dont have any card entries")
		return
	}
	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}

	// parse input
	pname := inputSelect("Card name: ", pnames)
	resp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{Pname: pname})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
	}

	// decrypting and converting to Card struct
	cleaned, err := helpers.FromPdata(resp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode card: %s\n", err.Error())
	}
	card := cleaned.(structs.Card)

	// print data
	upass_pretty, err := json.MarshalIndent(card, "", " ")
	if err != nil {
		log.Printf("ERROR cant encode upass JSON : %s\n", err.Error())
	} else {
		log.Println(string(upass_pretty))
	}
}
