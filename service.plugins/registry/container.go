package registry

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/domain/service"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/interface/persistence"
	"github.com/reecerussell/monzo-plus-plus/service.plugins/usecase"
)

// Service names.
const (
	ServicePluginUsecase = "plugin_usecase"
)

// Build builds a dependency injection container with the predefined services.
func Build() *di.Container {
	return di.New([]*di.Service{
		&di.Service{
			Name:     ServicePluginUsecase,
			Builder:  buildPluginUsecase,
			Lifetime: di.LifetimeSingleton,
		},
	}...)
}

func buildPluginUsecase(ctn *di.Container) (interface{}, error) {
	repo := persistence.NewPluginRepository()
	serv := service.NewPluginService()

	return usecase.NewPluginUsecase(serv, repo), nil
}
