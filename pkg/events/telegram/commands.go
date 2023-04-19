package telegram

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/url"
	"strings"
	"telegram-bot/solte.lab/pkg/clients/telegram"
	e "telegram-bot/solte.lab/pkg/errhandler"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/storage"
	"time"
)

const (
	CmdStart          = "/start"
	CmdHelp           = "/help"
	CmdRndWords       = "/words"
	CmdTopic          = "/topic"
	CmdSetTopic       = "/setTopic"
	CmdSetLanguage    = "/setLang"
	CmdPhraseOfTheDay = "/phraseOfDay"
)

func (p *Processor) doCmd(text string, user *models.User) error {
	p.logger.Info("Get new command", zap.String("content", text), zap.String("From User", user.Name))

	cmd, arg := parseCommand(text)
	//if isAddCmd(text) {
	//	return p.savePage(chatID, text, username)
	//}

	switch cmd {
	case CmdStart:
		return p.sendHello(user)
	case CmdHelp:
		return p.sendHelp(user)
	case CmdRndWords:
		return p.randomWords(user)
	case CmdPhraseOfTheDay:
		return p.phraseOfTheDay(user, arg)
	case CmdSetLanguage:
		return p.setLang(user, arg)
	default:
		return p.tg.SendMessage(user.ChatID, msgUnknownCommand)
	}
}

func (p *Processor) randomWords(user *models.User) (err error) {
	defer func() { err = e.WrapIfErr("can't execute command: random page", err) }()

	sendMsg := newMessageSender(user.ChatID, p.tg)

	word, err := p.storage.PickRandomWords(user)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}

	if errors.Is(err, storage.ErrNoSavedPages) {
		err = sendMsg(msgNoDataInStorage)
		if err != nil {
			return err
		}
		return nil
	}

	if err = sendMsg(word.Suomi); err != nil {
		return err
	}

	time.Sleep(3 * time.Second)

	//TODO more languages logic
	if user.IsLanguageEnglish() {
		if err = sendMsg(concatStrings(word.Suomi, word.English)); err != nil {
			return err
		}
		return nil
	}

	if err = sendMsg(concatStrings(word.Suomi, word.Russian)); err != nil {
		return err
	}

	return nil
}

func (p *Processor) phraseOfTheDay(user *models.User, arg string) (err error) {
	defer func() { err = e.WrapIfErr("can't execute command: for CmdPhraseOfTheDay", err) }()

	sendMsg := newMessageSender(user.ChatID, p.tg)
	if err = sendMsg(fmt.Sprintf("Phrase of the day\nMukavaa päivää - Hava a nice day")); err != nil {
		return err
	}

	return nil
}

func (p *Processor) setLang(user *models.User, arg string) (err error) {
	defer func() { err = e.WrapIfErr("can't execute command: set language", err) }()

	sendMsg := newMessageSender(user.ChatID, p.tg)

	if arg == "" {
		err = sendMsg(msgMissingArgument)
		if err != nil {
			return err
		}
		return nil
	}

	if err = user.CheckLanguage(arg); err != nil {
		err = sendMsg(msgUnsupportedLang)
		if err != nil {
			return err
		}

		return nil
	}

	err = p.storage.SetUserLanguage(user)
	if err != nil {
		return err
	}

	err = sendMsg(msgSettingApplied)
	if err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendHelp(user *models.User) error {
	return p.tg.SendMessage(user.ChatID, msgHelp)
}

func (p *Processor) sendHello(user *models.User) error {
	return p.tg.SendMessage(user.ChatID, msgHello)
}

func newMessageSender(ChatID int, tg *telegram.Client) func(string) error {
	return func(msg string) error {
		return tg.SendMessage(ChatID, msg)
	}
}

func parseCommand(text string) (cmd, arg string) {
	text = strings.TrimSpace(text)
	input := strings.Split(text, " ")

	if len(input) > 1 {
		return input[0], input[1]
	}

	return input[0], ""
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

//func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
//	defer func() { err = e.WrapIfErr("can't execute command: save page", err) }()
//
//	sendMsg := newMessageSender(chatID, p.tg)
//
//	page := storage.Page{
//		URL:      pageURL,
//		UserName: username,
//	}
//
//	isExist, err := p.storage.IsExist(&page)
//	if err != nil {
//		return err
//	}
//
//	if isExist {
//		return sendMsg(msgAlreadyExists)
//	}
//
//	if err := p.storage.Save(&page); err != nil {
//		return err
//	}
//
//	if err := sendMsg(msgPageSaved); err != nil {
//		return err
//	}
//
//	return nil
//}

//func (p *Processor) sendRandom(chatID int, username string) (err error) {
//	defer func() { err = e.WrapIfErr("can't execute command: random page", err) }()
//
//	sendMsg := newMessageSender(chatID, p.tg)
//
//	page, err := p.storage.PickRandom(username)
//	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
//		return err
//	}
//
//	if errors.Is(err, storage.ErrNoSavedPages) {
//		err := sendMsg(msgNoSavedPages)
//		if err != nil {
//			return err
//		}
//		return nil
//	}
//
//	if err := sendMsg(page.URL); err != nil {
//		return err
//	}
//
//	return p.storage.Remove(page)
//}
