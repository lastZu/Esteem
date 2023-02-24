package main

import (
	"flag"
	"github.com/lastZu/Esteem/internal/app/clients/telegram"
	"log"

	"github.com/lastZu/Esteem/internal/app/utils"

	"github.com/BurntSushi/toml"
)

var (
	configPath string
)

const (
	telegramBotHost = "api.telegram.org"
)

func init() {
	flag.StringVar(
		&configPath,
		"config-path",
		"configs/utils.toml",
		"path to config file")
}

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
	telegramClien := telegram.New(telegramBotHost, token)
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
