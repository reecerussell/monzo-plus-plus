package registry

import (
	"github.com/reecerussell/monzo-plus-plus/libraries/di"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/service"
	"github.com/reecerussell/monzo-plus-plus/service.auth/interface/persistence"
	"github.com/reecerussell/monzo-plus-plus/service.auth/password"
	"github.com/reecerussell/monzo-plus-plus/service.auth/usecase"
)

// Service names.
const (
	ServiceUserUsecase           = "user_usecase"
	ServicePermissionsRepository = "permissions_repository"
	ServicePasswordService       = "password_service"
	ServiceUserAuthUsecase       = "user_auth_usecase"
	ServiceRoleUsecase           = "role_usecase"
)

// Build builds the di container with the predefined services.
func Build() *di.Container {
	return di.New([]*di.Service{
		&di.Service{
			Name:     ServicePasswordService,
			Builder:  buildPasswordService,
			Lifetime: di.LifetimeSingleton,
		},
		&di.Service{
			Name:     ServicePermissionsRepository,
			Builder:  buildPermissionRepository,
			Lifetime: di.LifetimeSingleton,
		},
		&di.Service{
			Name:     ServiceUserAuthUsecase,
			Builder:  buildUserAuthUsecase,
			Lifetime: di.LifetimeSingleton,
		},
		&di.Service{
			Name:     ServiceUserUsecase,
			Builder:  buildUserUsecase,
			Lifetime: di.LifetimeSingleton,
		},
		&di.Service{
			Name:     ServiceRoleUsecase,
			Builder:  buildRoleUsecase,
			Lifetime: di.LifetimeSingleton,
		},
	}...)
}

func buildPasswordService(ctn *di.Container) (interface{}, error) {
	h := password.NewHasher()
	s := password.NewService(password.DefaultOptions, h)

	return s, nil
}

func buildPermissionRepository(ctn *di.Container) (interface{}, error) {
	return persistence.NewPermissionRepository(), nil
}

func buildUserAuthUsecase(ctn *di.Container) (interface{}, error) {
	ps := ctn.Resolve(ServicePasswordService).(password.Service)
	repo := persistence.NewUserRepository()
	serv := service.NewUserService()

	return usecase.NewUserAuthUsecase(ps, repo, serv)
}

func buildUserUsecase(ctn *di.Container) (interface{}, error) {
	repo := persistence.NewUserRepository()
	roles := persistence.NewRoleRepository()
	serv := service.NewUserService()
	ps := ctn.Resolve(ServicePasswordService).(password.Service)
	auth := ctn.Resolve(ServiceUserAuthUsecase).(usecase.UserAuthUsecase)

	return usecase.NewUserUsecase(repo, roles, serv, ps, auth), nil
}

func buildRoleUsecase(ctn *di.Container) (interface{}, error) {
	repo := persistence.NewRoleRepository()
	serv := service.NewRoleService()
	perms := persistence.NewPermissionRepository()

	return usecase.NewRoleUsecase(repo, serv, perms), nil
}
