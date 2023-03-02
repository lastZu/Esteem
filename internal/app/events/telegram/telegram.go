package telegram

import (
	"errors"
	"github.com/lastZu/Esteem/internal/app/clients/telegram"
	"github.com/lastZu/Esteem/internal/app/events"
	"github.com/lastZu/Esteem/internal/app/storage"
	"github.com/lastZu/Esteem/lib/e"
)

type Processor struct {
	client  *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	UserName string
}

var (
	ErrUnknownEventType = errors.New("unknown event type")
	ErrUnknownMetaType  = errors.New("unknown meta type")
)

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		client:  client,
		offset:  0,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.client.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	result := make([]events.Event, 0, len(updates))
	for _, update := range updates {
		result = append(result, event(update))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return result, nil
}

func (p *Processor) Process(event events.Event) error {
	if event.Type != events.Message {
		return e.Wrap("can't process message", ErrUnknownEventType)
	}

	p.processMessage(event)

	return nil
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("can't process message", err)
	}

	return nil
}

func processMessage(event events.Event) (Meta, error) {
	result, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("can't get meta", ErrUnknownMetaType)
	}

	return result, nil
}

func event(update telegram.Update) events.Event {
	updateType := fetchType(update)
	result := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}

	if updateType == events.Unknown {
		result.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			UserName: update.Message.From.UserName,
		}
	}

	return result

}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}

	return events.Message
}
