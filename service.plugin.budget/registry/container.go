package registry

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/interface/persistence"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/usecase"
)

// Service names.
const (
	PreferencesUsecaseService = "preferences_usecase"
)

func Build() *di.Container {
	return di.New([]*di.Service{
		&di.Service{
			Name:     PreferencesUsecaseService,
			Builder:  buildPreferencesUsecase,
			Lifetime: di.LifetimeSingleton,
		},
	}...)
}

func buildPreferencesUsecase(ctn *di.Container) (interface{}, error) {
	repo := persistence.NewPreferencesRepository()

	return usecase.NewPreferencesUsecase(repo), nil
}
