package client

import (
	"fmt"
	"log"
)

var commands = map[string]string{
	"register":  "create new user",
	"login":     "login to server",
	"help":      "list all available commands",
	"exit/quit": "exit application",
}

func help() {
	fmt.Println("AVAILABLE COMMANDS:")
	for cmd, desc := range commands {
		log.Printf("%s %s", cmd, desc)
	}
}
