package authorized

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"telegram-bot/solte.lab/pkg/api/tools"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/storage/emsql"
)

type Handlers struct {
	tools *tools.Tools
	view  *View
}

func New(db emsql.OPContract, fileServerPath string) *Handlers {
	view := NewView(fileServerPath, "index", "home.gohtml", "login.gohtml")
	t := tools.New(db)
	return &Handlers{
		tools: t,
		view:  view,
	}
}

func (h *Handlers) Register() (string, chi.Router) {
	routes := chi.NewRouter()
	routes.Use(h.tools.Middlewares.SetRequestID)
	routes.Use(h.tools.Middlewares.LogRequest)

	routes.Group(func(r chi.Router) {
		r.Use(h.tools.Middlewares.AuthenticateUser)
		r.Get("/test", h.login())
		r.Post("/add-words", h.handleAddWords())
	})

	routes.Group(func(r chi.Router) {
		r.Get("/", h.handleAdminPage())
	})

	filesDir := http.Dir(staticContent)
	h.fileServer(routes, "/static", filesDir)
	//routes.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("../../static"))))

	return "/admin", routes
}

func (h *Handlers) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h.tools.Responder.Send(w, r, http.StatusOK, "login successfully")
	}
}

func (h *Handlers) handleAddWords() http.HandlerFunc {
	type request struct {
		Input string `json:"words_string"`
	}
	p := models.NewParser()

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		var words []models.Words

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			h.tools.Responder.Error(w, r, http.StatusBadRequest, err)
			return
		}

		words, err := p.Parse(req.Input)
		if err != nil {
			h.tools.Responder.Error(w, r, http.StatusBadRequest, err)
			h.tools.Logger.Error("can't parse input data", zap.Error(err))
			return
		}

		if words == nil {
			h.tools.Responder.Error(w, r, http.StatusBadRequest, fmt.Errorf("bad input data"))
			h.tools.Logger.Error("empty words slice")
			return
		}

		if err := h.tools.Storage.Auth.AddNewWordsToDataBase(words); err != nil {
			h.tools.Responder.Error(w, r, http.StatusBadRequest, err)
			h.tools.Logger.Error("can't add word to storage", zap.Error(err))
			return
		}

		h.tools.Responder.Send(w, r, http.StatusOK, "Words added successfully")
	}
}

func (h *Handlers) fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rCtx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rCtx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func (h *Handlers) handleAdminPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("Session")
		if err != nil {
			err := h.view.renderView(w, r, nil)
			if err != nil {
				h.tools.Logger.Error("can't render view", zap.Error(err))
				h.tools.Responder.Error(w, r, http.StatusInternalServerError, err)
			}
			return
		}
		user, err := h.tools.Storage.Auth.FindBySessionToken(cookie.Value)
		if err != nil {
			err := h.view.renderView(w, r, nil)
			if err != nil {
				h.tools.Logger.Error("can't render view", zap.Error(err))
				h.tools.Responder.Error(w, r, http.StatusInternalServerError, err)
			}
			return
		}
		err = h.view.renderView(w, r, user.Name)
		if err != nil {
			h.tools.Logger.Error("can't render view", zap.Error(err))
			h.tools.Responder.Error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}
