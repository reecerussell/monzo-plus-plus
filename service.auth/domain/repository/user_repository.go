package repository

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
)

// UserRepository is used to manage persistent data in and from
// the MySQL database.
type UserRepository interface {
	Get(id string) (*model.User, errors.Error)
	GetByUsername(username string) (*model.User, errors.Error)
	GetList(term string) ([]*model.User, errors.Error)
	GetPending(term string) ([]*model.User, errors.Error)
	EnsureExists(id string) errors.Error
	Insert(u *model.User) errors.Error
	Update(u *model.User) errors.Error
	Delete(id string) errors.Error
}
