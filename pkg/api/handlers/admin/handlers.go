package admin

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
	"telegram-bot/solte.lab/pkg/api/tools"
	"telegram-bot/solte.lab/pkg/models"
	"telegram-bot/solte.lab/pkg/storage/emsql"
	"time"
)

type Admin struct {
	tools *tools.Tools
}

func New(db emsql.OPContract) *Admin {
	return &Admin{
		tools: tools.New(db),
	}
}

func (a *Admin) Register() (string, chi.Router) {
	routes := chi.NewRouter()
	routes.Use(a.tools.Middlewares.SetRequestID)
	routes.Use(a.tools.Middlewares.LogRequest)

	routes.Post("/login", a.handleUserLogin())
	routes.Post("/create_user", a.handleUsersCreate())
	routes.Post("/logout", a.handleUserLogOut)
	routes.Post("/test", a.handleTest())

	return "/api/authentication", routes
}

func (a *Admin) handleTest() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a.tools.Responder.Send(w, r, http.StatusOK, "test")
	}
}

func (a *Admin) handleUserLogOut(w http.ResponseWriter, r *http.Request) {
	cookie := http.Cookie{
		Name:     "Session",
		Value:    "",
		HttpOnly: true,
		Expires:  time.Now().Add(-time.Hour),
		Path:     "/",
	}
	http.SetCookie(w, &cookie)
	a.tools.Responder.Send(w, r, http.StatusOK, nil)
}

func (a *Admin) handleUserLogin() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := request{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			a.tools.Responder.Error(w, r, http.StatusBadRequest, err)
			return
		}
		user, err := a.tools.Storage.Auth.Authenticate(req.Email, req.Password)
		if err != nil {
			a.tools.Logger.Error("can't authenticate user", zap.Error(err))
			a.tools.Responder.Error(w, r, http.StatusUnauthorized, fmt.Errorf("incorrect email or password"))
			return
		}

		err = a.tools.Storage.Auth.SessionSave(user)
		if err != nil {
			a.tools.Responder.Error(w, r, http.StatusInternalServerError, fmt.Errorf("something goes wrong"))
			return
		}

		a.tools.Storage.Auth.Sanitize(user)

		cookie := http.Cookie{
			Name:     "Session",
			Value:    user.Token,
			HttpOnly: true,
			Expires:  time.Now().Add(4 * time.Hour),
			Path:     "/",
		}
		http.SetCookie(w, &cookie)
		data := map[string]string{
			"name":  user.Name,
			"email": user.Email,
		}
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		a.tools.Responder.Send(w, r, http.StatusOK, data)
	}
}

func (a *Admin) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			a.tools.Responder.Error(w, r, http.StatusBadRequest, err)
			return
		}

		user := &models.Admin{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
		}

		err := a.tools.Storage.Auth.CreateUser(user)
		if err != nil {
			a.tools.Responder.Error(w, r, http.StatusBadRequest, err)
			return
		}

		data := fmt.Sprintf("User %v has been created", user.Email)
		r.Header.Add("Content-Type", "application/json")
		a.tools.Responder.Send(w, r, http.StatusCreated, data)
	}
}
