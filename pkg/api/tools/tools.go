package tools

import (
	"go.uber.org/zap"
	"net/http"
	"telegram-bot/solte.lab/pkg/api/middleware"
	"telegram-bot/solte.lab/pkg/logging"
	"telegram-bot/solte.lab/pkg/services"
	"telegram-bot/solte.lab/pkg/storage/storagewrapper/postgresql"
)

type Responder interface {
	Send(w http.ResponseWriter, r *http.Request, code int, data interface{})
	Error(w http.ResponseWriter, r *http.Request, code int, err error)
}

type Tools struct {
	Logger      *zap.Logger
	Middlewares *middleware.Middlewares
	Responder   Responder
	Storage     *services.Services
}

func New(alias string) *Tools {
	st, _ := postgresql.GetStorage(alias)
	service := services.New(st)
	rs := NewResponder()
	return &Tools{
		Logger:      logging.GetLogger(),
		Middlewares: middleware.New(service, rs),
		Responder:   rs,
		Storage:     service,
	}
}
