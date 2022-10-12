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
		command := promptGetInput("command: ", notEmpty, false)
		fmt.Printf("mstorage: %v\n", mstorage)
		switch command {
		case "register":
			register(ctx, gclient)
		case "login":
			login(ctx, gclient, &mstorage)
		case "help":
			help()
		case "exit", "quit":
			os.Exit(0)
		default:
			log.Printf("'%s' is not supported", command)
			help()
		}

	}

}
