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

	// setting up gk for user
	ctx := context.Background()
	setup(mstorage, ctx, gclient)

	// infinite loop for interactive cli
	for {
		command := getInput("command: ", notEmpty, false)
		switch command {
		// Authentication
		case "register":
			err := register(ctx, gclient)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
			} else {
				fmt.Printf("register succsessful")
			}
		case "login":
			err := login(ctx, gclient, mstorage)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
			} else {
				fmt.Println("login succsessful")
			}

		// MasterKey
		case "key-generate":
			err := keyGenerate(mstorage)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
			} else {
				fmt.Println("key-generate was successful")
			}
		case "key-load":
			err := keyLoad("", mstorage)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err.Error())
			} else {
				fmt.Println("key-load succsessful")
			}
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
		case "card-delete":
			cardDelete(mstorage, ctx, gclient)

		// Private string
		case "pstring-add":
			pstringCreate(mstorage, ctx, gclient)
		case "pstring-get":
			pstringGet(mstorage, ctx, gclient)
		case "pstring-update":
			pstringUpdate(mstorage, ctx, gclient)
		case "pstring-delete":
			pstringDelete(mstorage, ctx, gclient)

		// Private file
		case "pfile-add":
			pfileAdd(mstorage, ctx, gclient)
		case "pfile-get":
			pfileGet(mstorage, ctx, gclient)
		case "pfile-update":
			pfileUpdate(mstorage, ctx, gclient)
		case "pfile-delete":
			pfileDelete(mstorage, ctx, gclient)

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

// setup setups gk for user: login/register and loading/generating master key
func setup(mstorage *structs.MemStorage, ctx context.Context, gclient *structs.Gclient) {
	log.Printf("Welcome to gophkeeper. Let`s set you up")

	// Register/Login
	answer := inputSelect("Do you want to register or log in?",
		[]string{"register", "login"})
	if answer == "register" {
		log.Println("Registering:")
		err := register(ctx, gclient)
		if err != nil {
			log.Fatalf("ERROR: %s", err.Error())
		}
		log.Panicln("Register succsessful.")
	}
	log.Println("Logging in:")
	err := login(ctx, gclient, mstorage)
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}
	log.Println("Login succsessful")

	// Load/Generate master key
	log.Println("Loading master key")
	err = keyLoad("", mstorage)
	if err != nil {
		log.Fatalf("ERROR: %s", err.Error())
	}

	log.Println("You are ready to go :)")
	log.Println("Print `help` to list available commands")
}
