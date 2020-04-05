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

type UserController struct {
	userUsecase     usecase.UserUsecase
	userAuthUsecase usecase.UserAuthUsecase
}

func NewUserController(ctn *di.Container, r *mux.Router) *UserController {
	uu := ctn.Resolve(registry.ServiceUserUsecase).(usecase.UserUsecase)
	uau := ctn.Resolve(registry.ServiceUserAuthUsecase).(usecase.UserAuthUsecase)

	c := &UserController{
		userUsecase:     uu,
		userAuthUsecase: uau,
	}

	r.HandleFunc("/users", c.HandleGetList).Methods("GET")
	r.HandleFunc("/users/pending", c.HandleGetPending).Methods("GET")
	r.HandleFunc("/users/{id}", c.HandleGet).Methods("GET")
	r.HandleFunc("/users", c.HandleCreate).Methods("POST")
	r.HandleFunc("/users/register", c.HandleRegister).Methods("POST")
	r.HandleFunc("/users", c.HandleUpdate).Methods("PUT")
	r.HandleFunc("/users/changepassword", c.HandleChangePassword).Methods("POST")
	r.HandleFunc("/users/enable/{id}", c.HandleEnable).Methods("POST")
	r.HandleFunc("/users/roles", c.HandleAddToRole).Methods("POST")
	r.HandleFunc("/users/roles", c.HandleRemoveFromRole).Methods("DELETE")
	r.HandleFunc("/users/roles/{id}", c.HandleGetRoles).Methods("GET")
	r.HandleFunc("/users/availableRoles/{id}", c.HandleGetAvailableRoles).Methods("GET")
	r.HandleFunc("/users/plugin", c.HandleEnablePlugin).Methods("POST")
	r.HandleFunc("/users/plugin", c.HandleDisablePlugin).Methods("DELETE")
	r.HandleFunc("/users/{id}", c.HandleDelete).Methods("DELETE")

	return c
}

func (c *UserController) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

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

	id := mux.Vars(r)["id"]

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

	id := mux.Vars(r)["id"]

	roles, err := c.userUsecase.GetRoles(r.Context(), id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&roles)
}

func (c *UserController) HandleGetAvailableRoles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

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

func (c *UserController) HandleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	ctx := r.Context()
	err := c.userUsecase.Delete(ctx, id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
