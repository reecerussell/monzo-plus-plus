package model

import (
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/datamodel"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
)

type Permission struct {
	id   int
	name string
}

func (p *Permission) GetID() int {
	return p.id
}

func (p *Permission) GetName() string {
	return p.name
}

func (p *Permission) DTO() *dto.Permission {
	return &dto.Permission{
		ID:   p.id,
		Name: p.name,
	}
}

func (p *Permission) DataModel() *datamodel.Permission {
	return &datamodel.Permission{
		ID:   p.id,
		Name: p.name,
	}
}

func PermissionFromDataModel(d *datamodel.Permission) *Permission {
	return &Permission{
		id:   d.ID,
		name: d.Name,
	}
}
