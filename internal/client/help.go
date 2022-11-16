package client

import (
	"fmt"
	"log"
)

var commands = map[string]string{
	"register":     "create new user",
	"login":        "login to server",
	"key-generate": "generate masterkey file",
	"key-load":     "load masterkey from file to memory",
	"key-print":    "get runtime masterkey information",
	"upass-add":    "generate user/password pair",
	"upass-get":    "retrive user/password pair",
	"help":         "list all available commands",
	"exit":         "exit application",
}

func help() {
	fmt.Println("AVAILABLE COMMANDS:")
	for cmd, desc := range commands {
		log.Printf("\t* %s - %s", cmd, desc)
	}
}
