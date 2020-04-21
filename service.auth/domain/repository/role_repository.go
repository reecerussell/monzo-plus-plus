package repository

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
)

type RoleRepository interface {
	Get(id string) (*model.Role, errors.Error)
	GetList(term string) ([]*model.Role, errors.Error)
	GetForUser(userID string) ([]*model.Role, errors.Error)
	GetAvailableForUser(userID string) ([]*model.Role, errors.Error)
	EnsureExists(id string) errors.Error
	Insert(r *model.Role) errors.Error
	Update(r *model.Role) errors.Error
	Delete(id string) errors.Error
}
