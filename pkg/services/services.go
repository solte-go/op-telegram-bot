package services

import (
	"telegram-bot/solte.lab/pkg/services/authentication"
	"telegram-bot/solte.lab/pkg/storage/emsql"
)

type Services struct {
	Auth authentication.Service
}

func New(st emsql.AdminsContract) *Services {
	return &Services{
		Auth: authentication.New(st),
	}
}
