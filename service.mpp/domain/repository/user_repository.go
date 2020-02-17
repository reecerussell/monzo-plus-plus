package repository

import (
	"errors"

	"github.com/reecerussell/monzo-plus-plus/service.mpp/domain/model"
)

// User data errors.
var (
	ErrUserNotFoundWithStateToken = errors.New("no user was found with this state token")
	ErrNoUserWasUpdated           = errors.New("no user was updated: user does not exist")
)

type UserRepository interface {
	FindByStateToken(token string) (*model.User, error)
	Get(userID string) (*model.User, error)
	Create(u *model.User) error
	Update(u *model.User) error
}
