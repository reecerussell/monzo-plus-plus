package model

import (
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/datamodel"
)

type Role struct {
	id   string
	name string
}

func (r *Role) DataModel() *datamodel.Role {
	return &datamodel.Role{
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
