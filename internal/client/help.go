package client

import (
	"fmt"
	"log"
)

var authCmds = map[string]string{
	"register": "create new user",
	"login":    "login to server",
}

var masterKeyCmds = map[string]string{
	"key-generate": "generate masterkey file",
	"key-load":     "load masterkey from file to memory",
	"key-print":    "get runtime masterkey information",
}

var upassCmds = map[string]string{
	"upass-add":    "generate user/password pair",
	"upass-get":    "retrive user/password pair",
	"upass-update": "update user/password pair",
	"upass-delete": "delete user/password pair",
}

var otherCmds = map[string]string{
	"help": "list all available commands",
	"exit": "exit application",
}

func help() {
	fmt.Println("Authentication commands:")
	for cmd, desc := range authCmds {
		log.Printf("\t* %s - %s", cmd, desc)
	}

	fmt.Println("Master key commands:")
	for cmd, desc := range masterKeyCmds {
		log.Printf("\t* %s - %s", cmd, desc)
	}

	fmt.Println("Upass commands")
	for cmd, desc := range upassCmds {
		log.Printf("\t* %s - %s", cmd, desc)
	}

	fmt.Println("Other commands:")
	for cmd, desc := range otherCmds {
		log.Printf("\t* %s - %s", cmd, desc)
	}

}
