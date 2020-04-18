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

type UserController struct {
	userUsecase     usecase.UserUsecase
	userAuthUsecase usecase.UserAuthUsecase
}

func NewUserController(ctn *di.Container, r *routing.Router) *UserController {
	uu := ctn.Resolve(registry.ServiceUserUsecase).(usecase.UserUsecase)
	uau := ctn.Resolve(registry.ServiceUserAuthUsecase).(usecase.UserAuthUsecase)

	c := &UserController{
		userUsecase:     uu,
		userAuthUsecase: uau,
	}

	r.GetFunc("/users", c.HandleGetList)
	r.GetFunc("/users/pending", c.HandleGetPending)
	r.GetFunc("/users/{id}", c.HandleGet)
	r.PostFunc("/users", c.HandleCreate)
	r.PostFunc("/users/register", c.HandleRegister)
	r.PutFunc("/users", c.HandleUpdate)
	r.PostFunc("/users/changepassword", c.HandleChangePassword)
	r.PostFunc("/users/enable/{id}", c.HandleEnable)
	r.PostFunc("/users/roles", c.HandleAddToRole)
	r.DeleteFunc("/users/roles", c.HandleRemoveFromRole)
	r.GetFunc("/users/roles/{id}", c.HandleGetRoles)
	r.GetFunc("/users/availableRoles/{id}", c.HandleGetAvailableRoles)
	r.PostFunc("/users/plugin", c.HandleEnablePlugin)
	r.DeleteFunc("/users/plugin", c.HandleDisablePlugin)
	r.PostFunc("/users/account", c.HandleSetAccount)
	r.GetFunc("/users/accounts/{id}", c.HandleGetAccounts)
	r.DeleteFunc("/users/{id}", c.HandleDelete)

	return c
}

func (c *UserController) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := routing.Vars(r)["id"]

	ctx := r.Context()
	user, err := c.userUsecase.Get(ctx, id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&user)
}

func (c *UserController) HandleGetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	term := r.URL.Query().Get("term")

	ctx := r.Context()
	users, err := c.userUsecase.GetList(ctx, term)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&users)
}

func (c *UserController) HandleGetPending(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	term := r.URL.Query().Get("term")

	ctx := r.Context()
	users, err := c.userUsecase.GetPending(ctx, term)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&users)
}

func (c *UserController) HandleCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var d dto.CreateUser
	_ = json.NewDecoder(r.Body).Decode(&d)

	ctx := r.Context()
	user, err := c.userUsecase.Create(ctx, &d)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&user)
}

func (c *UserController) HandleRegister(w http.ResponseWriter, r *http.Request) {
	var d dto.CreateUser
	_ = json.NewDecoder(r.Body).Decode(&d)

	stateToken, err := c.userAuthUsecase.Register(&d)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(stateToken))
}

func (c *UserController) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var d dto.UpdateUser
	_ = json.NewDecoder(r.Body).Decode(&d)

	ctx := r.Context()
	err := c.userUsecase.Update(ctx, &d)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *UserController) HandleChangePassword(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var d dto.ChangePassword
	_ = json.NewDecoder(r.Body).Decode(&d)

	ctx := r.Context()
	err := c.userUsecase.ChangePassword(ctx, &d)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *UserController) HandleEnable(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := routing.Vars(r)["id"]

	ctx := r.Context()
	err := c.userUsecase.Enable(ctx, id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *UserController) HandleAddToRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ur dto.UserRole
	_ = json.NewDecoder(r.Body).Decode(&ur)

	err := c.userUsecase.AddToRole(r.Context(), &ur)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
	}
}

func (c *UserController) HandleRemoveFromRole(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ur dto.UserRole
	_ = json.NewDecoder(r.Body).Decode(&ur)

	err := c.userUsecase.RemoveFromRole(r.Context(), &ur)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
	}
}

func (c *UserController) HandleGetRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := routing.Vars(r)["id"]

	roles, err := c.userUsecase.GetRoles(r.Context(), id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&roles)
}

func (c *UserController) HandleGetAvailableRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := routing.Vars(r)["id"]

	roles, err := c.userUsecase.GetAvailableRoles(r.Context(), id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&roles)
}

func (c *UserController) HandleEnablePlugin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ur dto.UserPlugin
	_ = json.NewDecoder(r.Body).Decode(&ur)

	err := c.userUsecase.EnablePlugin(r.Context(), &ur)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
	}
}

func (c *UserController) HandleDisablePlugin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ur dto.UserPlugin
	_ = json.NewDecoder(r.Body).Decode(&ur)

	err := c.userUsecase.DisablePlugin(r.Context(), &ur)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
	}
}

func (c *UserController) HandleGetAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := routing.Vars(r)["id"]

	accounts, err := c.userUsecase.GetAccounts(r.Context(), id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&accounts)
}

func (c *UserController) HandleSetAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var ua dto.UserAccount
	_ = json.NewDecoder(r.Body).Decode(&ua)

	err := c.userUsecase.SetAccount(r.Context(), &ua)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
	}
}

func (c *UserController) HandleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := routing.Vars(r)["id"]

	ctx := r.Context()
	err := c.userUsecase.Delete(ctx, id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
