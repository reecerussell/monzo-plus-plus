package repository

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
)

type PermissionsRepository interface {
	LoadCollections() map[int][]string
	Get(id int) (*model.Permission, errors.Error)
}
