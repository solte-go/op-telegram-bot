package services

import (
	"telegram-bot/solte.lab/pkg/services/authentication"
	"telegram-bot/solte.lab/pkg/storage"
)

type Services struct {
	Auth authentication.Service
}

func New(st storage.Administrators) *Services {
	return &Services{
		Auth: authentication.New(st),
	}
}
