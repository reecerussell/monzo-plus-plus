package registry

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/di"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/interface/persistence"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/usecase"
)

// Service names.
const (
	UserUsecaseService = "user_usecase_service"
)

func Build() *di.Container {
	return di.New([]*di.Service{
		&di.Service{
			Name:     UserUsecaseService,
			Builder:  buildUserUsecase,
			Lifetime: di.LifetimeSingleton,
		},
	}...)
}

func buildUserUsecase(ctn *di.Container) (interface{}, error) {
	repo := persistence.NewUserRepository().(repository.UserRepository)
	usecase := usecase.NewUserUsecase(repo)

	return usecase, nil
}
