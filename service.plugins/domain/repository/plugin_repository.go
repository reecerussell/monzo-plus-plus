package repository

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/model"
)

// PluginRepository is used to read and write data to the MySQL database.
type PluginRepository interface {
	GetList(term string) ([]*model.Plugin, errors.Error)
	Get(id string) ([]*model.Plugin, errors.Error)
	Create(p *model.Plugin) errors.Error
	Update(p *model.Plugin) errors.Error
	Delete(id string) errors.Error
}
