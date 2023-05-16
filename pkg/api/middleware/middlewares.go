package middleware

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"net/http"
	"telegram-bot/solte.lab/pkg/logging"
	"telegram-bot/solte.lab/pkg/services"
	"time"
)

type ctxKey int8

const (
	CtxKeyUser ctxKey = iota
	ctxKeyRequestID
	CtxErrorKey
)

type responseWriter struct {
	http.ResponseWriter
	code int
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type Responder interface {
	Send(w http.ResponseWriter, r *http.Request, code int, data interface{})
	Error(w http.ResponseWriter, r *http.Request, code int, err error)
}

type Middlewares struct {
	*zap.Logger
	*services.Services
	Responder
}

func New(st *services.Services, Respond Responder) *Middlewares {
	return &Middlewares{
		logging.GetLogger(),
		st,
		Respond,
	}
}

func (m *Middlewares) SetRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (m *Middlewares) LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := m.Logger.With(
			zap.String("type", "api_request"),
			zap.Any("src_addr", r.RemoteAddr),
			zap.Any("request_id", r.Context().Value(ctxKeyRequestID)),
		)

		logger.Info(fmt.Sprintf("started %s %s", r.Method, r.RequestURI))

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		level := zap.DebugLevel
		switch {
		case rw.code >= 500:
			level = zap.ErrorLevel
		case rw.code == 400 || rw.code >= 402:
			level = zap.WarnLevel
		default:
		}
		logger.Log(
			level,
			fmt.Sprintf("completed with %d %s in %v", rw.code, http.StatusText(rw.code), time.Now().Sub(start)),
		)
	})
}

func (m *Middlewares) AuthenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Session")
		if err != nil {
			//next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxErrorKey, err)))
			m.Responder.Error(w, r, http.StatusUnauthorized, err)
			return
		}

		user, err := m.Auth.FindBySessionToken(cookie.Value)
		if err != nil {
			//next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxErrorKey, err)))
			m.Responder.Error(w, r, http.StatusUnauthorized, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, user.Email)))
	})
}
