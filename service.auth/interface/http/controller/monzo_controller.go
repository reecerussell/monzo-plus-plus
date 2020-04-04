package controller

import (
	"net/http"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"

	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"

	"github.com/reecerussell/monzo-plus-plus/service.auth/registry"
	"github.com/reecerussell/monzo-plus-plus/service.auth/usecase"

	"github.com/gorilla/mux"
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
)

type MonzoController struct {
	userUsecase usecase.UserUsecase
}

func NewMonzoController(ctn *di.Container, r *mux.Router) *MonzoController {
	userUsecase := ctn.Resolve(registry.ServiceUserUsecase).(usecase.UserUsecase)

	c := &MonzoController{
		userUsecase: userUsecase,
	}

	r.HandleFunc("/monzo/login", c.HandleLogin).Methods("GET")
	r.HandleFunc("/monzo/callback", c.HandleCallback).Methods("GET")

	return c
}

func (c *MonzoController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	state := r.URL.Query().Get("state")
	if state == "" {
		http.NotFound(w, r)
		return
	}

	monzo.Login(w, r, state)
}

func (c *MonzoController) HandleCallback(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	code := vals.Get("code")
	state := vals.Get("state")

	if code == "" || state == "" {
		http.NotFound(w, r)
		return
	}

	err := c.userUsecase.Login(code, state)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	http.Redirect(w, r, "http://localhost:3000/#login", http.StatusPermanentRedirect)
}
