package storage

import (
	"crypto/sha1"
	"fmt"
	"github.com/lastZu/Esteem/lib/e"
	"io"
)

type Storage interface {
	Save(page *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(page *Page) error
	IsExists(page *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
}

func (p Page) Hash() (hash string, err error) {
	defer func() { err = e.WrapIfErr("can't calculate hash", err) }()

	hashSum := sha1.New()

	if _, err := io.WriteString(hashSum, p.URL); err != nil {
		return "", err
	}
	if _, err := io.WriteString(hashSum, p.UserName); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hashSum.Sum(nil)), nil
}
