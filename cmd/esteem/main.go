package main

import (
	"flag"
	"log"

	"github.com/lastZu/Esteem/internal/app/clients/telegram"
	event_consumer "github.com/lastZu/Esteem/internal/app/consumer/event-consumer"
	tEvent "github.com/lastZu/Esteem/internal/app/events/telegram"
	"github.com/lastZu/Esteem/internal/app/storage/files"
)

var (
	configPath string
)

const (
	telegramBotHost = "api.telegram.org"
	storagePath     = "storage"
	batchSize       = 100
)

func init() {
	flag.StringVar(
		&configPath,
		"config-path",
		"configs/utils.toml",
		"path to config file",
	)
}

func main() {
	token := mustToken()
	telegramClien := telegram.New(telegramBotHost, token)

	eventsProcessor := tEvent.New(telegramClien, files.New(storagePath))

	log.Print("service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stoped", err)
	}
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
