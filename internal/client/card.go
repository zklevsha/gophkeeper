package client

import (
	"context"
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
