package usecase

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
)

// UserRoleUsecase is an interface providing methods to manage
// a user's role assignments.
type UserRoleUsecase interface {
	AddToRole(d *dto.UserRole) errors.Error
	RemoveFromRole(d *dto.UserRole)
}
