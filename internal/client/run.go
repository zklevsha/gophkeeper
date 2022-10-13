package client

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/zklevsha/gophkeeper/internal/structs"
)

// Run starts client interactive promt
func Run(gclient structs.Gclient) {
	fmt.Println("Welcome to gophkeeper")
	fmt.Println("Enter 'help' to get list of available commands")
	ctx := context.Background()
	mstorage := structs.MemStorage{}
	for {
		command := getInput("command: ", notEmpty, false)
		switch command {
		case "register":
			register(ctx, gclient)
		case "login":
			login(ctx, gclient, &mstorage)
		case "key-generate":
			keyGenerate(&mstorage)
		case "key-load":
			keyLoad("", &mstorage)
		case "key-print":
			keyPrint(&mstorage)
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
