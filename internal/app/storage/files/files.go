package files

import (
	"encoding/gob"
	"errors"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/lastZu/Esteem/internal/app/storage"
	"github.com/lastZu/Esteem/lib/e"
)

type Storage struct {
	basePath string
}

const (
	defaultPermission = 0774
)

func New(basePath string) Storage {
	return Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page", err) }()

	filePath := filepath.Join(s.basePath, page.UserName)

	if err := os.Mkdir(filePath, defaultPermission); err != nil {
		return err
	}

	fileName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return err
	}

	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page", err) }()

	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(len(files))
	file := files[n]

	filePath := filepath.Join(path, file.Name())
	return s.decodePage(filePath)
}

func (s Storage) Remove(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't remove file", err) }()

	fileName, err := fileName(page)
	if err != nil {
		return err
	}

	path := filepath.Join(s.basePath, page.UserName, fileName)
	if err := os.Remove(path); err != nil {
		return err
	}

	return nil
}

func (s Storage) IsExists(page *storage.Page) (exist bool, err error) {
	defer func() { err = e.WrapIfErr("can't check file", err) }()

	fileName, err := fileName(page)
	if err != nil {
		return false, err
	}

	path := filepath.Join(s.basePath, page.UserName, fileName)
	switch _, err = os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, err
	}

	return true, nil
}

func fileName(page *storage.Page) (string, error) {
	return page.Hash()
}

func (s Storage) decodePage(filePath string) (p *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't decode page", err) }()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var page storage.Page
	if err := gob.NewDecoder(file).Decode(&page); err != nil {
		return nil, err
	}

	return &page, nil
}
