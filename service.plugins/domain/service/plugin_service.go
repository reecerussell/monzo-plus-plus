package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/model"
)

type PluginService struct {
	db *sql.DB
}

func NewPluginService() *PluginService {
	return new(PluginService)
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

	ctx := context.Background()
	stmt, err := ps.db.PrepareContext(ctx, query)
	if err != nil {
		return errors.InternalError(err)
	}
	defer stmt.Close()

	var count int

	err = stmt.QueryRowContext(ctx, args...).Scan(&count)
	if err != nil {
		return errors.InternalError(err)
	}

	if count > 0 {
		return errors.BadRequest(fmt.Sprintf("the name '%s' is already taken", p.GetName()))
	}

	return nil
}
