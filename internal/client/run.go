package client

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/zklevsha/gophkeeper/internal/structs"
)

func notEmpty(input string) error {
	if len(input) <= 0 {
		return errors.New("input is empty")
	}
	return nil
}

func Run(gclient structs.Gclient) {
	ctx := context.Background()
	for {
		command := promptGetInput("command: ", notEmpty, false)
		switch command {
		case "login":
			login(ctx, gclient)
		case "exit", "quit":
			os.Exit(0)
		default:
			log.Printf("'%s' is not supported", command)
		}

	}

}
