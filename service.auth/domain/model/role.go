package model

import (
	"github.com/google/uuid"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/datamodel"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
)

type Role struct {
	id   string
	name string
}

func NewRole(d *dto.CreateRole) (*Role, errors.Error) {
	id, _ := uuid.NewRandom()
	r := new(Role)

	r.id = id.String()

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

func RoleFromDataModel(dm *datamodel.Role) *Role {
	return &Role{
		id:   dm.ID,
		name: dm.Name,
	}
}
