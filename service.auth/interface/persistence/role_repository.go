package persistence

import (
	"database/sql"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
)

type roleRepository struct {
	db *sql.DB
}

func NewRoleRepository() repository.RoleRepository {
	return new(roleRepository)
}

func (rr *roleRepository) Get(id string) (*model.Role, errors.Error) {
	return nil, nil
}

func (rr *roleRepository) GetList(id string) ([]*model.Role, errors.Error) {
	return nil, nil
}

func (rr *roleRepository) EnsureExists(id string) errors.Error {
	return nil
}

func (rr *roleRepository) Insert(r *model.Role) errors.Error {
	return nil
}

func (rr *roleRepository) Update(r *model.Role) errors.Error {
	return nil
}

func (rr *roleRepository) Delete(id string) errors.Error {
	return nil
}
