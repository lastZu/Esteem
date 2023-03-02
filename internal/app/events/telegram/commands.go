package telegram

import (
	"net/url"
	"strings"
)

type data struct {
	text     string
	chatID   int
	userName string
}

type command func(data2 data) error

var (
	commands = map[string]command{
		"add": add,
	}
)

func add(data2 data) error {

}

func (p *Processor) doCmd(text string, chatID int, userName string) error {
	text = strings.TrimSpace(text)

	// Log

	if isAddCmd(text) {
		err := commands[text](data{})
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
