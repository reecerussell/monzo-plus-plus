package controller

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/registry"
	"github.com/reecerussell/monzo-plus-plus/service.auth/usecase"
)

type TokenController struct {
	uau usecase.UserAuthUsecase
}

func NewTokenController(ctn *di.Container, r *mux.Router) *TokenController {
	uau := ctn.Resolve(registry.ServiceUserAuthUsecase).(usecase.UserAuthUsecase)

	c := &TokenController{
		uau: uau,
	}

	r.HandleFunc("/token", c.HandleToken).Methods("POST")

	return c
}

// HandleToken handles HTTP POST requests to generate a JSON-Web token, for the
// given credentials.
//
// TODO: create a token struct, to include expiry dates etc.
func (c *TokenController) HandleToken(w http.ResponseWriter, r *http.Request) {
	var cred dto.UserCredential
	_ = json.NewDecoder(r.Body).Decode(&cred)

	token, err := c.uau.GenerateToken(&cred)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&token)
}
