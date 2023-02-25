package telegram

import "github.com/lastZu/Esteem/internal/app/clients/telegram"

type Processor struct {
	client *telegram.Client
	offset int
	// storage
}
