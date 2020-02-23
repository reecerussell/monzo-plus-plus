package repository

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
)

// UserRoleRepository creates and delete links between users and roles.
type UserRoleRepository interface {
	Insert(userID, roleID string) errors.Error
	Delete(userID, roleID string) errors.Error
}
