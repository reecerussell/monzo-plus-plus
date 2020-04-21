package controller

import (
	"net/http"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"
	"github.com/reecerussell/monzo-plus-plus/libraries/routing"
	"github.com/reecerussell/monzo-plus-plus/service.auth/registry"
	"github.com/reecerussell/monzo-plus-plus/service.auth/usecase"
)

type MonzoController struct {
	userAuthUsecase usecase.UserAuthUsecase
}

func NewMonzoController(ctn *di.Container, r *routing.Router) *MonzoController {
	userAuthUsecase := ctn.Resolve(registry.ServiceUserAuthUsecase).(usecase.UserAuthUsecase)

	c := &MonzoController{
		userAuthUsecase: userAuthUsecase,
	}

	r.GetFunc("/monzo/login", c.HandleLogin)
	r.GetFunc("/monzo/callback", c.HandleCallback)

	return c
}

func (c *MonzoController) HandleLogin(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	state, id := query.Get("state"), query.Get("id")
	if state == "" && id == "" {
		http.NotFound(w, r)
		return
	}

	if id != "" {
		var err errors.Error
		state, err = c.userAuthUsecase.GetStateToken(id)
		if err != nil {
			errors.HandleHTTPError(w, r, err)
			return
		}
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

	err := c.userAuthUsecase.Login(code, state)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	http.Redirect(w, r, "http://localhost:3000/#login", http.StatusPermanentRedirect)
}
