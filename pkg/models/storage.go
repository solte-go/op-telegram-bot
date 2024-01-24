package models

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	e "telegram-bot/solte.lab/pkg/errhandler"
)

// epäjärjestelmällistyttämättömyydellänsäkäänköhän
var regex = regexp.MustCompile(`^([a-zA-Z]{1,60},){2}([а-яА-Я]{1,60},){1}([a-zA-Z]{1,60})`)

var (
	ErrNotEnoughArguments = errors.New("not enough number of arguments")
	ErrEmptyData          = errors.New("string shouldn't be empty")
	ErrBadDataInLine      = errors.New("bad data in line")
)

type Page struct {
	UserID   int
	UserName string
	URLId    int
	URL      string
}

type Words struct {
	Topic   string
	Letter  string
	Suomi   string
	Russian string
	English string
}

type lineRules func() error

func runLineRulesFunc(fns ...lineRules) error {
	for _, fn := range fns {
		if err := fn(); err != nil {
			return err
		}
	}
	return nil
}

type Parser struct {
	line  string
	words []string
}

func NewParser() Parser {
	return Parser{}
}

func (parser *Parser) Parse(input string) ([]Words, error) {
	lines := strings.Split(input, ";")

	words := make([]Words, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}

		parser.line = line

		if err := runLineRulesFunc(
			parser.isEmpty,
			parser.normalize,
			parser.lineSyntastic,
			parser.isEnoughArgs,
			parser.isArgEmpty,
		); err != nil {
			return nil, e.Wrap("can't parse line", err)
		}

		newWord, err := parser.populate()
		if err != nil {
			return nil, e.Wrap("can't populate word", err)
		}

		words = append(words, newWord)
	}

	return words, nil
}

func (parser *Parser) isEmpty() error {
	if parser.line == "" {
		return ErrEmptyData
	}
	return nil
}

func (parser *Parser) lineSyntastic() error {
	if !regex.MatchString(parser.line) {
		return e.Wrap(parser.line, ErrBadDataInLine)
	}
	//p.line = strings.Trim(p.line, ";")
	return nil
}

func (parser *Parser) normalize() error {
	parser.line = strings.TrimSpace(
		strings.ToLower(
			strings.ReplaceAll(parser.line, " ", "")))
	return nil
}

func (parser *Parser) isEnoughArgs() error {
	parser.words = strings.Split(parser.line, ",")

	if len(parser.words)-1 < 3 {
		return ErrNotEnoughArguments
	}
	return nil
}

func (parser *Parser) isArgEmpty() error {
	for _, word := range parser.words {
		if word == "" {
			return ErrEmptyData
		}
	}
	return nil
}

func (parser *Parser) populate() (Words, error) {
	w := Words{}
	w.Letter = strings.TrimSpace(parser.words[0][0:1])
	w.Suomi = strings.TrimSpace(parser.words[0])
	w.English = strings.TrimSpace(parser.words[1])
	w.Russian = strings.TrimSpace(parser.words[2])
	w.Topic = strings.TrimSpace(parser.words[3])
	return w, nil
}

func (p *Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("can't create hash", err)
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("can't create hash", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
