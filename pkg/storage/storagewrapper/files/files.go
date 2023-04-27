package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/storage"
)

type Storage struct {
	basePath string
}

const defaultPermissions = 0o774

func New(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s *Storage) Save(page *storage.Page) (err error) {
	defer func() { err = e.WrapIfErr("can't save page to file", err) }()

	filePath := filepath.Join(s.basePath, page.UserName)
	if err := os.MkdirAll(filePath, defaultPermissions); err != nil {
		return err
	}

	fName, err := fileName(page)
	if err != nil {
		return err
	}

	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	err = gob.NewEncoder(file).Encode(page)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = e.WrapIfErr("can't pick random page from file", err) }()

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

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s *Storage) Remove(page *storage.Page) error {
	fName, err := fileName(page)
	if err != nil {
		return e.Wrap("can't remove file", err)
	}

	path := filepath.Join(s.basePath, page.UserName, fName)

	if err := os.Remove(path); err != nil {
		return e.Wrap(fmt.Sprintf("can't remove file: %s", path), err)
	}
	return nil
}

func (s *Storage) IsExist(page *storage.Page) (bool, error) {
	fName, err := fileName(page)
	if err != nil {
		return false, e.Wrap("can't find file", err)
	}

	path := filepath.Join(s.basePath, page.UserName, fName)

	switch _, err = os.Stat(path); {
	case errors.Is(os.ErrNotExist, err):
		return false, nil
	case err != nil:
		return false, e.Wrap(fmt.Sprintf("can't find file: %s", path), err)
	}
	return true, nil
}

func (s *Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, e.Wrap("cant' decode file", err)
	}
	defer func() { _ = file.Close() }()

	p := &storage.Page{}

	if err := gob.NewDecoder(file).Decode(p); err != nil {
		return nil, e.Wrap("cant' decode file", err)
	}
	return p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
