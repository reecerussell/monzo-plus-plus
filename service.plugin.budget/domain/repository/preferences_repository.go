package repository

import (
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/domain/model"
)

// PreferencesRepository is used to manage persistent data for
// the preferences domain in the budget plugin.
type PreferencesRepository interface {
	Get(userID string) (*model.Preferences, error)
	Update(p *model.Preferences) error
}
