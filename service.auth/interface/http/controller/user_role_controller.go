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

type UserRoleController struct {
	u usecase.UserRoleUsecase
}

func NewUserRoleController(ctn *di.Container, r *mux.Router) *UserRoleController {
	u := ctn.Resolve(registry.ServiceUserRoleUsecase).(usecase.UserRoleUsecase)

	c := &UserRoleController{
		u: u,
	}

	r.HandleFunc("/user/role", c.HandleAdd).Methods("POST")
	r.HandleFunc("/user/role", c.HandleRemove).Methods("POST")

	return c
}

func (c *UserRoleController) HandleAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var d dto.UserRole
	_ = json.NewDecoder(r.Body).Decode(&d)

	ctx := r.Context()
	err := c.u.AddToRole(ctx, &d)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *UserRoleController) HandleRemove(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var d dto.UserRole
	_ = json.NewDecoder(r.Body).Decode(&d)

	ctx := r.Context()
	err := c.u.RemoveFromRole(ctx, &d)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
