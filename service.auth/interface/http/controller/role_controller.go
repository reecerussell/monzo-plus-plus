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

type RoleController struct {
	u usecase.RoleUsecase
}

func NewRoleController(ctn *di.Container, r *mux.Router) *RoleController {
	u := ctn.Resolve(registry.ServiceRoleUsecase).(usecase.RoleUsecase)

	c := &RoleController{
		u: u,
	}

	r.HandleFunc("/roles", c.HandleGetList).Methods("GET")
	r.HandleFunc("/roles/{id}", c.HandleGet).Methods("GET")
	r.HandleFunc("/roles", c.HandleCreate).Methods("POST")
	r.HandleFunc("/roles", c.HandleUpdate).Methods("PUT")
	r.HandleFunc("/roles/permission", c.HandleAddPermission).Methods("POST")
	r.HandleFunc("/roles/permission", c.HandleRemovePermission).Methods("DELETE")
	r.HandleFunc("/roles/{id}", c.HandleDelete).Methods("DELETE")

	return c
}

func (c *RoleController) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	ctx := r.Context()
	role, err := c.u.Get(ctx, id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&role)
}

func (c *RoleController) HandleGetList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	term := r.URL.Query().Get("term")

	ctx := r.Context()
	roles, err := c.u.GetList(ctx, term)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&roles)
}

func (c *RoleController) HandleCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var d dto.CreateRole
	_ = json.NewDecoder(r.Body).Decode(&d)

	ctx := r.Context()
	role, err := c.u.Create(ctx, &d)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	json.NewEncoder(w).Encode(&role)
}

func (c *RoleController) HandleUpdate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var d dto.Role
	_ = json.NewDecoder(r.Body).Decode(&d)

	ctx := r.Context()
	err := c.u.Update(ctx, &d)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *RoleController) HandleAddPermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rp dto.RolePermission
	_ = json.NewDecoder(r.Body).Decode(&rp)

	err := c.u.AddPermission(r.Context(), &rp)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *RoleController) HandleRemovePermission(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var rp dto.RolePermission
	_ = json.NewDecoder(r.Body).Decode(&rp)

	err := c.u.RemovePermission(r.Context(), &rp)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *RoleController) HandleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := mux.Vars(r)["id"]

	ctx := r.Context()
	err := c.u.Delete(ctx, id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
