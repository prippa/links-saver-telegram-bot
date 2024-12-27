package telegram

import (
	"errors"
	"fmt"
	"links-saver-telegram-bot/storage"
	"log"
	"net/url"
	"os"
	"strings"
)

const (
	RandCmd  = "/rand"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("get new command: '%s' from user: %s", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RandCmd:
		return p.randCmd(chatID, username)
	case HelpCmd:
		return p.helpCmd(chatID)
	case StartCmd:
		return p.startCmd(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCmd)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return fmt.Errorf("error checking page existence: %w", err)
	}
	if isExists {
		return p.tg.SendMessage(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return fmt.Errorf("error saving page: %w", err)
	}

	if err := p.tg.SendMessage(chatID, msgSaved); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return nil
}

func (p *Processor) randCmd(chatID int, username string) error {
	page, err := p.storage.PickRandom(username)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) || errors.Is(err, storage.ErrNoSavedPages) {
			return p.tg.SendMessage(chatID, msgNoSavedPages)
		}

		return fmt.Errorf("error picking random page: %w", err)
	}

	if page == nil {
		return p.tg.SendMessage(chatID, msgNoSavedPages)
	}

	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}

	return p.storage.Remove(page)
}

func (p *Processor) helpCmd(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) startCmd(chatID int) error {
	return p.tg.SendMessage(chatID, msgStart)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
