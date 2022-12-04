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
	cvc := getInput("CVV/CVC(XXX):", isCardCVC, false)
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

// cardGet retreives Credit card entry from the server
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

	// Getting pdata from server
	pname := inputSelect("Card name: ", pnames)
	pdataID := entries[pname]
	resp, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{PdataID: pdataID})
	if err != nil {
		log.Printf("ERROR: cant retrive pdata from server: %s\n", err.Error())
		return
	}

	// decrypting and converting to Card struct
	cleaned, err := helpers.FromPdata(resp.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode card: %s\n", err.Error())
		return
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

// cardUpdate Credit card entry from the server
func cardUpdate(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {

	err := reqCheck(mstorage)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// getting list of existing card entries
	entries, err := listPnames(ctx, gclient, "card")
	if err != nil {
		log.Printf("ERROR: cant retrive list of existing card entries: %s", err.Error())
	}
	if len(entries) == 0 {
		log.Printf("You dont have any card entries")
		return
	}

	// getting selected card
	var pnames []string
	for pname := range entries {
		pnames = append(pnames, pname)
	}
	pname := inputSelect("Card name: ", pnames)
	pdataID := entries[pname]
	getPdataResponse, err := gclient.Pdata.GetPdata(ctx, &pb.GetPdataRequest{PdataID: pdataID})
	if err != nil {
		log.Printf("ERROR: cant get card entry: %s\n", err.Error())
		return
	}
	cleaned, err := helpers.FromPdata(getPdataResponse.Pdata, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant decode card: %s\n", err.Error())
		return
	}
	oldCard := cleaned.(structs.Card)

	// getting new data from input
	var newCard structs.Card
	// name
	name := getInput(fmt.Sprintf("name[%s]:", oldCard.Name), any, false)
	if name == "" {
		newCard.Name = oldCard.Name
	} else {
		newCard.Name = name
	}
	// card number
	number := getInput(fmt.Sprintf("number[%s]:", oldCard.Number), isCardNumberOrEmpty, false)
	if number == "" {
		newCard.Number = oldCard.Number
	} else {
		newCard.Number = number
	}
	// card holder
	holder := getInput(fmt.Sprintf("holder(%s):", oldCard.Holder), isCardHolderOrEmpty, false)
	if holder == "" {
		newCard.Holder = oldCard.Holder
	} else {
		newCard.Holder = holder
	}
	// card expiration date
	expire := getInput(fmt.Sprintf("expire[%s]:", oldCard.Expire), isCardExpireOrEmpty, false)
	if expire == "" {
		newCard.Expire = oldCard.Expire
	} else {
		newCard.Expire = expire
	}
	// card CVV/CVC number
	cvc := getInput(fmt.Sprintf("CVV/CVC[%s]:)", oldCard.CVC), isCardCVCorEmpty, false)
	if cvc == "" {
		newCard.CVC = oldCard.CVC
	} else {
		newCard.CVC = cvc
	}
	// tags
	tagsJson, err := json.Marshal(oldCard.Tags)
	if err != nil {
		log.Printf("ERROR: cant parse old tags: %s\n", err.Error())
		return
	}
	tagsStr := getInput(fmt.Sprintf("new tags[%s]", tagsJson), isTags, false)
	var tagsNew map[string]string
	if tagsStr == "" {
		tagsNew = oldCard.Tags
	} else {
		tagsNew, err = getTags(tagsStr)
		if err != nil {
			log.Printf("ERROR: cant convert tags %s\n", err.Error())
			return
		}
	}
	newCard.Tags = tagsNew

	// converting new card to pdata
	pdata, err := helpers.ToPdata("card", newCard, mstorage.MasterKey)
	if err != nil {
		log.Printf("ERROR: cant convert to pdata %s", err.Error())
	}
	pdata.ID = pdataID

	// sending data to server
	resp, err := gclient.Pdata.UpdatePdata(ctx, &pb.UpdatePdataRequest{Pdata: pdata})
	if err != nil {
		log.Printf("ERROR: cant send message to server: %s\n", err.Error())
		return
	}
	log.Println(resp.Response)

}
