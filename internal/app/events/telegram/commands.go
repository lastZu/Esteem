package telegram

import (
	"errors"
	"net/url"
	"strings"

	"github.com/lastZu/Esteem/internal/app/storage"
	"github.com/lastZu/Esteem/lib/e"
)

type data struct {
	text     string
	chatID   int
	userName string
}

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)

	// Log

	if isAddCmd(text) {
		p.savePage(text, chatID, userName)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, userName)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.client.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(pageURL string, chatID int, userName string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: userName,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.client.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.client.SendMessage(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.client.SendMessage(chatID, msgNoSavedPage)
	}

	if err := p.client.SendMessage(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) (err error) {
	return p.client.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) (err error) {
	return p.client.SendMessage(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
