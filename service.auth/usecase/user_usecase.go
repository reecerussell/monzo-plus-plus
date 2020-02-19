package usecase

import (
	"context"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
)

// UserUsecase is a high-level interface providing methods to perform
// CRUD operations, as well as, more specific operations on the User domain.
type UserUsecase interface {
	Create(d *dto.CreateUser) (*dto.User, errors.Error)
	Get(id string) (*dto.User, errors.Error)
	GetList(term string) ([]*dto.User, errors.Error)
	GetPending(term string) ([]*dto.User, errors.Error)
	Update(d *dto.UpdateUser) errors.Error
	Enable(id string) errors.Error
	Delete(id string) errors.Error
	WithUser(ctx context.Context, accessToken string) (context.Context, errors.Error)
}
