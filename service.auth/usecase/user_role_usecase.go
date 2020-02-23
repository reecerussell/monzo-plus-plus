package usecase

import (
	"context"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.auth/permission"
)

// UserRoleUsecase is an interface providing methods to manage
// a user's role assignments.
type UserRoleUsecase interface {
	AddToRole(ctx context.Context, d *dto.UserRole) errors.Error
	RemoveFromRole(ctx context.Context, d *dto.UserRole) errors.Error
}

type userRoleUsecase struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	repo     repository.UserRoleRepository
}

func NewUserRoleUsecase(userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	repo repository.UserRoleRepository) UserRoleUsecase {
	return &userRoleUsecase{
		userRepo: userRepo,
		roleRepo: roleRepo,
		repo:     repo,
	}
}

// AddToRole creates a link from the given user to the role. Both the
// user's and the role's existance are checked.
//
// An error is returned if either the user executing the operation has
// insufficient permissions, the user or role don't exist, or if there
// was an error inserting the record into the database.
func (uru *userRoleUsecase) AddToRole(ctx context.Context, d *dto.UserRole) errors.Error {
	if !permission.Has(ctx, permission.PermissionRoleManager) {
		return errors.Forbidden()
	}

	g := &errors.Group{}

	g.Go(func() errors.Error {
		return uru.userRepo.EnsureExists(d.UserID)
	})
	g.Go(func() errors.Error {
		return uru.roleRepo.EnsureExists(d.RoleID)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	err := uru.repo.Insert(d.UserID, d.RoleID)
	if err != nil {
		return err
	}

	return nil
}

// RemoveFromRole remove the link from the given user to the role.
// Both the user's and the role's existance are checked.
//
// An error is returned if either the user executing the operation has
// insufficient permissions, the user or role don't exist, or if there
// was an error inserting the record into the database.
func (uru *userRoleUsecase) RemoveFromRole(ctx context.Context, d *dto.UserRole) errors.Error {
	if !permission.Has(ctx, permission.PermissionRoleManager) {
		return errors.Forbidden()
	}

	g := &errors.Group{}

	g.Go(func() errors.Error {
		return uru.userRepo.EnsureExists(d.UserID)
	})
	g.Go(func() errors.Error {
		return uru.roleRepo.EnsureExists(d.RoleID)
	})

	if err := g.Wait(); err != nil {
		return err
	}

	err := uru.repo.Delete(d.UserID, d.RoleID)
	if err != nil {
		return err
	}

	return nil
}
