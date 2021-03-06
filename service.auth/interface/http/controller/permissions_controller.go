package controller

import (
	"net/http"

	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/routing"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.auth/permission"
	"github.com/reecerussell/monzo-plus-plus/service.auth/registry"
)

type PermissionsController struct {
	repo repository.PermissionsRepository
}

func NewPermissionsController(ctn *di.Container, r *routing.Router) *PermissionsController {
	repo := ctn.Resolve(registry.ServicePermissionsRepository).(repository.PermissionsRepository)

	c := &PermissionsController{repo}

	r.PostFunc("/permissions/flush", c.HandleFlush)

	return c
}

func (c *PermissionsController) HandleFlush(w http.ResponseWriter, r *http.Request) {
	if !permission.Has(r.Context(), permission.PermissionRoleManager) {
		errors.HandleHTTPError(w, r, errors.Forbidden())
	} else {
		permission.Build(c.repo)
	}
}
