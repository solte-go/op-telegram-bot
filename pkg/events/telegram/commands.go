package telegram

import (
	"errors"
	"net/url"
	"sort"
	"strings"
	"telegram-bot/solte.lab/pkg/clients/telegram"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/storage"
	"telegram-bot/solte.lab/pkg/storage/dialect"
	"time"

	"go.uber.org/zap"

	e "telegram-bot/solte.lab/pkg/errhandler"
)

const (
	CmdStart          = "/start"
	CmdHelp           = "/help"
	CmdHelpRu         = "/helpRu"
	CmdHelpEn         = "/helpEn"
	CmdRndWords       = "/words"
	CmdTopics         = "/topics"
	CmdSetTopic       = "/setTopic"
	CmdSetLanguage    = "/setLang"
	CmdPhraseOfTheDay = "/phraseOfDay"
)

func (r *Responder) doCmd(user *models.User) error {
	r.logger.Debug("get new command", zap.String("content", user.Cmd), zap.String("From User", user.Name))

	cmd, arg := parseCommand(user.Cmd)

	switch cmd {
	case CmdStart:
		return r.sendHello(user)
	case CmdHelp:
		return r.sendHelp(user)
	case CmdHelpRu:
		return r.sendHelpRu(user)
	case CmdHelpEn:
		return r.sendHelpEn(user)
	case CmdRndWords:
		return r.randomWords(user)
	case CmdPhraseOfTheDay:
		return r.phraseOfTheDay(user, arg)
	case CmdTopics:
		return r.sendTopics(user)
	case CmdSetLanguage:
		return r.setLang(user, arg)
	case CmdSetTopic:
		return r.setTopic(user, arg)
	default:
		return r.tg.SendMessage(user.ChatID, msgUnknownCommand)
	}
}

func (r *Responder) randomWords(user *models.User) (err error) {
	defer func() { err = e.WrapIfErr("can't execute command: random page", err) }()

	var word *models.Words
	sendMsg := newMessageSender(user.ChatID, r.tg)

	word, err = r.worker.PickRandomWord(user)
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

	// TODO more languages logic
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

func (r *Responder) phraseOfTheDay(user *models.User, arg string) (err error) {
	defer func() { err = e.WrapIfErr("can't execute command: for CmdPhraseOfTheDay", err) }()

	sendMsg := newMessageSender(user.ChatID, r.tg)
	if err = sendMsg("Phrase of the day\nMukavaa päivää - Hava a nice day"); err != nil {
		return err
	}

	return nil
}

func (r *Responder) setLang(user *models.User, arg string) (err error) {
	defer func() { err = e.WrapIfErr("can't execute command: set language", err) }()

	sendMsg := newMessageSender(user.ChatID, r.tg)

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

	err = r.worker.SetUserLanguage(user)
	if err != nil {
		return err
	}

	err = sendMsg(msgSettingApplied)
	if err != nil {
		return err
	}

	return nil
}

func (r *Responder) sendTopics(user *models.User) (err error) {
	defer func() { err = e.WrapIfErr("can't execute command: send topics", err) }()

	sendMsg := newMessageSender(user.ChatID, r.tg)

	topics, err := r.worker.GetTopics()
	if err != nil {
		return err
	}

	if len(topics) == 0 {
		err = sendMsg(msgNoDataInStorage)
		if err != nil {
			return err
		}
		return nil
	}

	sort.Strings(topics)

	if err = sendMsg(concatStringsAsList(topics...)); err != nil {
		return err
	}

	return nil
}

func (r *Responder) setTopic(user *models.User, arg string) (err error) {
	defer func() { err = e.WrapIfErr("can't execute command: set topic", err) }()

	sendMsg := newMessageSender(user.ChatID, r.tg)

	err = r.worker.SetUserTopic(user, arg)
	if err != nil && errors.Is(err, dialect.ErrUnsupportedTopic) {
		err = sendMsg(msgUnsupportedTopic)
		if err != nil {
			return err
		}
		return nil
	}

	if err != nil {
		return err
	}

	err = sendMsg(msgSettingApplied)
	if err != nil {
		return err
	}

	return nil
}

func (r *Responder) sendHelp(user *models.User) error {
	return r.tg.SendMessage(user.ChatID, msgHelp)
}

func (r *Responder) sendHelpRu(user *models.User) error {
	return r.tg.SendMessage(user.ChatID, msgHelpRu)
}

func (r *Responder) sendHelpEn(user *models.User) error {
	return r.tg.SendMessage(user.ChatID, msgHelpEn)
}

func (r *Responder) sendHello(user *models.User) error {
	return r.tg.SendMessage(user.ChatID, msgHello)
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
