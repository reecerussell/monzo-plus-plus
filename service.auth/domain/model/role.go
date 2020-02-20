package model

import (
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/datamodel"
)

type Role struct {
	id   string
	name string
}

func RoleFromDataModel(dm *datamodel.Role) *Role {
	return &Role{
		id:   dm.ID,
		name: dm.Name,
	}
}
