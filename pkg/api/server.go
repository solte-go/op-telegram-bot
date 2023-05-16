package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

type Server struct {
	router *chi.Mux
	logger *zap.Logger
}

type Handler interface {
	Register() (string, chi.Router)
}

func New(logger *zap.Logger) *Server {
	r := chi.NewRouter()
	return &Server{router: r, logger: logger}
}

func (s *Server) Run(ctx context.Context, port int, handlers ...Handler) {
	for _, h := range handlers {
		path, router := h.Register()
		s.router.Mount(path, router)
	}

	s.logger.Debug(fmt.Sprintf("Server running on port %d", port))
	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.router,
	}

	go func() {
		// graceful shutdown
		<-ctx.Done()
		s.logger.Warn("server shutdown")
		srv.Shutdown(context.Background())
	}()

	err := srv.ListenAndServe()
	if err != nil {
		s.logger.Error("server error", zap.Error(err))
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
