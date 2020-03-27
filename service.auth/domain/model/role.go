package model

import (
	"fmt"

	"github.com/google/uuid"

	"github.com/reecerussell/monzo-plus-plus/libraries/domain"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/datamodel"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/event"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/handler"
)

func init() {
	domain.RegisterEventHandler(&event.AddPermissionToRole{}, &handler.AddPermissionToRole{})
	domain.RegisterEventHandler(&event.RemovePermissionFromRole{}, &handler.RemovePermissionFromRole{})
}

type Role struct {
	domain.Aggregate

	id   string
	name string

	permissions []*Permission
}

func NewRole(d *dto.CreateRole) (*Role, errors.Error) {
	id, _ := uuid.NewRandom()
	r := new(Role)

	r.id = id.String()
	r.permissions = []*Permission{}

	err := r.UpdateName(d.Name)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (r *Role) Update(d *dto.Role) errors.Error {
	err := r.UpdateName(d.Name)
	if err != nil {
		return err
	}

	return nil
}

func (r *Role) UpdateName(name string) errors.Error {
	l := len(name)
	if l < 1 {
		return errors.BadRequest("name cannot be empty")
	}

	if l > 25 {
		return errors.BadRequest("name cannot be greater than 25 characters")
	}

	r.name = name

	return nil
}

func (r *Role) AddPermission(p *Permission) errors.Error {
	for _, ep := range r.permissions {
		if ep.GetID() == p.GetID() {
			return errors.BadRequest(fmt.Sprintf("permission '%s' has already been added", p.GetName()))
		}
	}

	r.RaiseEvent(&event.AddPermissionToRole{
		RoleID:       r.id,
		PermissionID: p.GetID(),
	})

	return nil
}

func (r *Role) RemovePermission(p *Permission) errors.Error {
	for _, ep := range r.permissions {
		if ep.GetID() == p.GetID() {
			r.RaiseEvent(&event.RemovePermissionFromRole{
				RoleID:       r.id,
				PermissionID: p.GetID(),
			})
			return nil
		}
	}

	return errors.BadRequest(fmt.Sprintf("permission '%s' has not already been added", p.GetName()))
}

func (r *Role) DataModel() *datamodel.Role {
	return &datamodel.Role{
		ID:   r.id,
		Name: r.name,
	}
}

func (r *Role) DTO() *dto.Role {
	return &dto.Role{
		ID:   r.id,
		Name: r.name,
	}
}

func RoleFromDataModel(dm *datamodel.Role, permissions ...*datamodel.Permission) *Role {
	perms := make([]*Permission, len(permissions))
	for i, p := range permissions {
		perms[i] = PermissionFromDataModel(p)
	}

	return &Role{
		id:          dm.ID,
		name:        dm.Name,
		permissions: perms,
	}
}
