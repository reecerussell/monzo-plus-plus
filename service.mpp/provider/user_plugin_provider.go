package provider

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/database"
	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
)

// UserPluginProvider is used to read
type UserPluginProvider struct {
	db *database.DB
}

// NewUserPluginProvider returns a new instance of UserPluginProvider.
func NewUserPluginProvider() *UserPluginProvider {
	return &UserPluginProvider{
		db: database.New(),
	}
}

// GetUserPlugins returns a list of plugin ids for a specific user.
func (p *UserPluginProvider) GetUserPlugins(userID string) ([]string, errors.Error) {
	query := "CALL get_user_plugins(?);"

	arr, err := p.db.Read(query, func(s database.ScannerFunc) (interface{}, errors.Error) {
		var id string

		err := s(&id)
		if err != nil {
			return nil, errors.InternalError(err)
		}

		return id, nil
	}, userID)

	if err != nil {
		return nil, err
	}

	ids := make([]string, len(arr))

	for i, id := range arr {
		ids[i] = id.(string)
	}

	return ids, nil
}
