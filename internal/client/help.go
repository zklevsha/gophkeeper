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
	"card-add":    "add credit card entry",
	"card-get":    "retrive card entry",
	"card-update": "update card entry",
	"card-delete": "delete card entry",
}

var pstringCmds = map[string]string{
	"pstring-add":    "add private string entry",
	"pstring-get":    "retrive private string entry",
	"pstring-update": "update pstring  entry",
	"pstring-delete": "delete pstring entry",
}

var pfileCmds = map[string]string{
	"pfile-add":    "add private file entry",
	"pfile-get":    "get private file entry",
	"pfile-update": "update private file entry",
	"pfile-delete": "delete private file entry",
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

	fmt.Println("Private string commands")
	for cmd, desc := range pstringCmds {
		log.Printf("\t* %s - %s", cmd, desc)
	}

	fmt.Println("Private file commands")
	for cmd, desc := range pfileCmds {
		log.Printf("\t* %s - %s", cmd, desc)
	}

	fmt.Println("Other commands:")
	for cmd, desc := range otherCmds {
		log.Printf("\t* %s - %s", cmd, desc)
	}

}
