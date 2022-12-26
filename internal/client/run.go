package client

import (
	"context"
	"fmt"
	"log"
	"os"
)

// Run starts client interactive promt
func Run(gclient *Gclient, mstorage *MemStorage) {

	// setting up gk for user
	ctx := context.Background()
	setup(ctx, mstorage, gclient)

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
			upassCreate(ctx, mstorage, gclient)
		case "upass-get":
			upassGet(ctx, mstorage, gclient)
		case "upass-update":
			upassUpdate(ctx, mstorage, gclient)
		case "upass-delete":
			upassDelete(ctx, mstorage, gclient)

		// Credit card
		case "card-add":
			cardCreate(ctx, mstorage, gclient)
		case "card-get":
			cardGet(ctx, mstorage, gclient)
		case "card-update":
			cardUpdate(ctx, mstorage, gclient)
		case "card-delete":
			cardDelete(ctx, mstorage,  gclient)

		// Private string
		case "pstring-add":
			pstringCreate(ctx, mstorage, gclient)
		case "pstring-get":
			pstringGet(ctx, mstorage, gclient)
		case "pstring-update":
			pstringUpdate(ctx, mstorage, gclient)
		case "pstring-delete":
			pstringDelete(ctx, mstorage, gclient)

		// Private file
		case "pfile-add":
			pfileAdd(ctx, mstorage, gclient)
		case "pfile-get":
			pfileGet( ctx, mstorage, gclient)
		case "pfile-update":
			pfileUpdate(ctx, mstorage, gclient)
		case "pfile-delete":
			pfileDelete(ctx, mstorage, gclient)

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
func setup(ctx context.Context, mstorage *MemStorage, gclient *Gclient) {
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
		log.Println("Register succsessful.")
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
