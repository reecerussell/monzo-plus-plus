package controller

import (
	"encoding/json"
	"net/http"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/routing"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/registry"
	"github.com/reecerussell/monzo-plus-plus/service.auth/usecase"
)

type TokenController struct {
	uau usecase.UserAuthUsecase
}

func NewTokenController(ctn *di.Container, r *routing.Router) *TokenController {
	uau := ctn.Resolve(registry.ServiceUserAuthUsecase).(usecase.UserAuthUsecase)

	c := &TokenController{
		uau: uau,
	}

	r.PostFunc("/token", c.HandleToken)
	r.GetFunc("/refresh", c.HandleRefresh)

	return c
}

// HandleToken handles HTTP POST requests to generate a JSON-Web token, for the
// given credentials.
func (c *TokenController) HandleToken(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var cred dto.UserCredential
	_ = json.NewDecoder(r.Body).Decode(&cred)

	token, err := c.uau.GenerateToken(&cred)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&token)
}

// HandleRefresh is used to refresh an access token.
func (c *TokenController) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	token, err := c.uau.RefreshToken(r.Context())
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&token)
}
