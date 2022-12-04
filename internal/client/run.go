package client

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zklevsha/gophkeeper/internal/structs"
)

// Run starts client interactive promt
func Run(gclient *structs.Gclient, mstorage *structs.MemStorage) {

	fmt.Println("Welcome to gophkeeper")
	fmt.Println("Enter 'help' to get list of available commands")
	ctx := context.Background()
	for {
		command := getInput("command: ", notEmpty, false)
		switch command {
		// Authentication
		case "register":
			register(ctx, gclient)
		case "login":
			login(ctx, gclient, mstorage)
		// MasterKey
		case "key-generate":
			keyGenerate(mstorage)
		case "key-load":
			keyLoad("", mstorage)
		case "key-print":
			keyPrint(mstorage)
		// Upass
		case "upass-add":
			upassCreate(mstorage, ctx, gclient)
		case "upass-get":
			upassGet(mstorage, ctx, gclient)
		case "upass-update":
			upassUpdate(mstorage, ctx, gclient)
		case "upass-delete":
			upassDelete(mstorage, ctx, gclient)
		// Credit card
		case "card-add":
			cardCreate(mstorage, ctx, gclient)
		case "card-get":
			cardGet(mstorage, ctx, gclient)
		case "card-update":
			cardUpdate(mstorage, ctx, gclient)
			// Other
		case "help":
			help()
		case "exit", "quit":
			os.Exit(0)
		default:
			log.Printf("command '%s' is not supported", command)
			help()
		}

	}

}
