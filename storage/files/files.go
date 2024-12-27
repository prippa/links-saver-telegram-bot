package files

import (
	"encoding/gob"
	"errors"
	"fmt"
	"links-saver-telegram-bot/storage"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

type Storage struct {
	basePath string
}

const defaultPerm = 0754

func New(basePath string) *Storage {
	return &Storage{
		basePath: basePath,
	}
}

func (s Storage) Save(page *storage.Page) (err error) {
	fPath := filepath.Join(s.basePath, page.UserName)

	if err := os.Mkdir(fPath, defaultPerm); err != nil && !os.IsExist(err) {
		return fmt.Errorf("error creating directory: %w", err)
	}

	fName, err := fuleName(page)
	if err != nil {
		return fmt.Errorf("error getting file name: %w", err)
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return fmt.Errorf("error encoding page: %w", err)
	}

	return nil
}

func (s Storage) PickRandom(userName string) (*storage.Page, error) {
	path := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	rand.New(rand.NewSource(time.Now().UnixNano()))
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func (s Storage) Remove(page *storage.Page) error {
	fName, err := fuleName(page)
	if err != nil {
		return fmt.Errorf("error getting file name: %w", err)
	}

	fPath := filepath.Join(s.basePath, page.UserName, fName)

	if err := os.Remove(fPath); err != nil {
		return fmt.Errorf("error removing file (%s) : %w", fPath, err)
	}

	return nil
}

func (s Storage) IsExists(page *storage.Page) (bool, error) {
	fName, err := fuleName(page)
	if err != nil {
		return false, fmt.Errorf("error getting file name: %w", err)
	}

	fPath := filepath.Join(s.basePath, page.UserName, fName)

	if _, err := os.Stat(fPath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, fmt.Errorf("error stating file (%s) : %w", fPath, err)
	}

	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var p storage.Page
	if err := gob.NewDecoder(file).Decode(&p); err != nil {
		return nil, fmt.Errorf("error decoding page: %w", err)
	}

	return &p, nil
}

func fuleName(p *storage.Page) (string, error) {
	return p.Hash()
}
