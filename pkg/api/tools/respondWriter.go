package tools

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
	"telegram-bot/solte.lab/pkg/logging"
)

type ResponseTemplates struct {
	logger *zap.Logger
}

func NewResponder() *ResponseTemplates {
	return &ResponseTemplates{
		logger: logging.GetLogger(),
	}
}

func (rt *ResponseTemplates) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	rt.Send(w, r, code, map[string]string{"error": err.Error()})
}

func (rt *ResponseTemplates) Send(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			rt.logger.Error("can't encode data", zap.Error(err))
			return
		}
	}
}
