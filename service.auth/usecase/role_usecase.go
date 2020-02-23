package usecase

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
)

type RoleUsecase interface {
	Get(id string) (*dto.Role, errors.Error)
	GetList(term string) ([]*dto.Role, errors.Error)
}
