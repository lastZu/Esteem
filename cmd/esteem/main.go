package main

import (
	"flag"
	"log"
)

func main() {
	token := mustToken()
}

func mustToken() string {
	token := flag.String(
		"t",
		"",
		"token for access to telegram bot",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("wrong token")
	}

	return *token
}
