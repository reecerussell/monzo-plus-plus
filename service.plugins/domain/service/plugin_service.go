package service

import (
	"fmt"

	"github.com/reecerussell/monzo-plus-plus/libraries/database"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/model"
)

type PluginService struct {
	db *database.DB
}

func NewPluginService() *PluginService {
	return &PluginService{
		db: database.New(),
	}
}

// EnsureUniqueName searches the database to ensure the given plugin
// has a unique name. A nil-error is returned is the plugin's name is
// unique, however, an error is returned if either there is a problem
// speaking to the database or the name is invalid.
func (ps *PluginService) EnsureUniqueName(p *model.Plugin) errors.Error {
	query := "SELECT COUNT(*) FROM plugins WHERE name LIKE ? AND id != ?;"
	args := []interface{}{
		p.GetName(),
		p.GetID(),
	}

	exists, err := ps.db.Exists(query, args...)
	if err != nil {
		return err
	}

	if exists {
		return errors.BadRequest(fmt.Sprintf("the name '%s' is already taken", p.GetName()))
	}

	return nil
}
