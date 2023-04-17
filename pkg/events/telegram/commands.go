package telegram

import (
	"errors"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"telegram-bot/solte.lab/pkg/clients/telegram"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/storage"
)

const (
	CmdStart = "/start"
	CmdHelp  = "/help"
	CmdRnd   = "/rnd"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)

	p.logger.Info("Get new command", zap.String("content", text), zap.String("From User", username))

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case CmdStart:
		return p.sendHello(chatID)
	case CmdRnd:
		return p.sendRandom(chatID, username)
	case CmdHelp:
		return p.sendHelp(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't execute command: save page", err) }()

	sendMsg := newMessageSender(chatID, p.tg)

	page := storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExist, err := p.storage.IsExist(&page)
	if err != nil {
		return err
	}

	if isExist {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(&page); err != nil {
		return err
	}

	if err := sendMsg(msgPageSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.WrapIfErr("can't execute command: random page", err) }()

	sendMsg := newMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		err := sendMsg(msgNoSavedPages)
		if err != nil {
			return err
		}
		return nil
	}

	if err := sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

func newMessageSender(ChatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(ChatID, msg)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
