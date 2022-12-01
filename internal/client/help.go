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
	"upass-add":    "add user/password entry",
	"upass-get":    "retrive user/password entry",
	"upass-update": "update user/password entry",
	"upass-delete": "delete user/password entry",
}

var cardCmds = map[string]string{
	"card-add": "add credit card entry",
	"card-get": "retrive card entry",
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

	fmt.Println("Credit card commands")
	for cmd, desc := range cardCmds {
		log.Printf("\t* %s - %s", cmd, desc)
	}

	fmt.Println("Other commands:")
	for cmd, desc := range otherCmds {
		log.Printf("\t* %s - %s", cmd, desc)
	}

}
