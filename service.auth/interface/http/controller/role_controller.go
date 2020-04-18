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

type RoleController struct {
	u usecase.RoleUsecase
}

func NewRoleController(ctn *di.Container, r *routing.Router) *RoleController {
	u := ctn.Resolve(registry.ServiceRoleUsecase).(usecase.RoleUsecase)

	c := &RoleController{
		u: u,
	}

	r.GetFunc("/roles", c.HandleGetList)
	r.GetFunc("/roles/{id}", c.HandleGet)
	r.PostFunc("/roles", c.HandleCreate)
	r.PutFunc("/roles", c.HandleUpdate)
	r.PostFunc("/roles/permission", c.HandleAddPermission)
	r.DeleteFunc("/roles/permission", c.HandleRemovePermission)
	r.GetFunc("/roles/permissions/{id}", c.HandleGetPermissions)
	r.GetFunc("/roles/availablePermissions/{id}", c.HandleGetAvailablePermissions)
	r.DeleteFunc("/roles/{id}", c.HandleDelete)

	return c
}

func (c *RoleController) HandleGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := routing.Vars(r)["id"]

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

func (c *RoleController) HandleGetPermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := routing.Vars(r)["id"]

	perms, err := c.u.GetPermissions(r.Context(), id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
	} else {
		json.NewEncoder(w).Encode(&perms)
	}
}

func (c *RoleController) HandleGetAvailablePermissions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := routing.Vars(r)["id"]

	perms, err := c.u.GetAvailablePermissions(r.Context(), id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
	} else {
		json.NewEncoder(w).Encode(&perms)
	}
}

func (c *RoleController) HandleDelete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	id := routing.Vars(r)["id"]

	ctx := r.Context()
	err := c.u.Delete(ctx, id)
	if err != nil {
		errors.HandleHTTPError(w, r, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
