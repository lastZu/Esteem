package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/lastZu/Esteem/internal/app/utils"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

func main() {
	flag.Parse()

	config := utils.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		log.Fatal(err)
	}

	engineer := utils.New(config)
	if err := engineer.Start(); err != nil {
		log.Fatal(err)
	}

	token := mustToken()

	fmt.Println(engineer)
	fmt.Println(token)
}

func init() {
	flag.StringVar(
		&configPath,
		"config-path",
		"configs/utils.toml",
		"path to config file")
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
