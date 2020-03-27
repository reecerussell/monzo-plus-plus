package usecase

import (
	"context"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/util"

	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/service"
	"github.com/reecerussell/monzo-plus-plus/service.auth/password"
	"github.com/reecerussell/monzo-plus-plus/service.auth/permission"
)

// UserUsecase is a high-level interface providing methods to perform
// CRUD operations, as well as, more specific operations on the User domain.
type UserUsecase interface {
	Create(ctx context.Context, d *dto.CreateUser) (*dto.User, errors.Error)
	Get(ctx context.Context, id string) (*dto.User, errors.Error)
	GetList(ctx context.Context, term string) ([]*dto.User, errors.Error)
	GetPending(ctx context.Context, term string) ([]*dto.User, errors.Error)
	Update(ctx context.Context, d *dto.UpdateUser) errors.Error
	ChangePassword(ctx context.Context, d *dto.ChangePassword) errors.Error
	Enable(ctx context.Context, id string) errors.Error
	Delete(ctx context.Context, id string) errors.Error
}

type userUsecase struct {
	repo repository.UserRepository
	serv *service.UserService
	ps   password.Service
}

// NewUserUsecase instantiates a new instance of UserUsecase with the given dependencies.
func NewUserUsecase(repo repository.UserRepository, serv *service.UserService, ps password.Service) UserUsecase {
	return &userUsecase{
		repo: repo,
		serv: serv,
		ps:   ps,
	}
}

func (uu *userUsecase) Create(ctx context.Context, d *dto.CreateUser) (*dto.User, errors.Error) {
	u, err := model.NewUser(d, uu.ps)
	if err != nil {
		return nil, err
	}

	err = uu.serv.ValidateUsername(u)
	if err != nil {
		return nil, err
	}

	err = uu.repo.Insert(u)
	if err != nil {
		return nil, err
	}

	return u.DTO(), nil
}

func (uu *userUsecase) Get(ctx context.Context, id string) (*dto.User, errors.Error) {
	currentUserID := ctx.Value(util.ContextKey("user_id"))
	if id != currentUserID && !permission.Has(ctx, permission.PermissionGetUser) {
		return nil, errors.Forbidden()
	}

	u, err := uu.repo.Get(id)
	if err != nil {
		return nil, err
	}

	return u.DTO(), nil
}

func (uu *userUsecase) GetList(ctx context.Context, term string) ([]*dto.User, errors.Error) {
	if !permission.Has(ctx, permission.PermissionGetUserList) {
		return nil, errors.Forbidden()
	}

	users, err := uu.repo.GetList(term)
	if err != nil {
		return nil, err
	}

	return convertToDTOs(users), nil
}

func (uu *userUsecase) GetPending(ctx context.Context, term string) ([]*dto.User, errors.Error) {
	if !permission.Has(ctx, permission.PermissionGetPendingUsers) {
		return nil, errors.Forbidden()
	}

	users, err := uu.repo.GetPending(term)
	if err != nil {
		return nil, err
	}

	return convertToDTOs(users), nil
}

func convertToDTOs(users []*model.User) []*dto.User {
	dtos := make([]*dto.User, len(users))

	for i, u := range users {
		dtos[i] = u.DTO()
	}

	return dtos
}

func (uu *userUsecase) Update(ctx context.Context, d *dto.UpdateUser) errors.Error {
	currentUserID := ctx.Value(util.ContextKey("user_id"))
	if d.ID != currentUserID && !permission.Has(ctx, permission.PermissionUpdateUser) {
		return errors.Forbidden()
	}

	u, err := uu.repo.Get(d.ID)
	if err != nil {
		return err
	}

	err = u.Update(d)
	if err != nil {
		return err
	}

	err = uu.serv.ValidateUsername(u)
	if err != nil {
		return err
	}

	err = uu.repo.Update(u)
	if err != nil {
		return err
	}

	return nil
}

// ChangePassword is used to change the password for the requests user. This method
// expected the context to be propogated with the logged in user, through the HTTP
// authentication middleware.
func (uu *userUsecase) ChangePassword(ctx context.Context, d *dto.ChangePassword) errors.Error {
	u := ctx.Value(util.ContextKey("user")).(*model.User)

	err := u.UpdatePassword(d.NewPassword, d.CurrentPassword, uu.ps)
	if err != nil {
		return err
	}

	err = uu.repo.Update(u)
	if err != nil {
		return err
	}

	return nil
}

func (uu *userUsecase) Enable(ctx context.Context, id string) errors.Error {
	if !permission.Has(ctx, permission.PermissionEnableUser) {
		return errors.Forbidden()
	}

	u, err := uu.repo.Get(id)
	if err != nil {
		return err
	}

	err = u.Enable()
	if err != nil {
		return err
	}

	err = uu.repo.Update(u)
	if err != nil {
		return err
	}

	return nil
}

func (uu *userUsecase) Delete(ctx context.Context, id string) errors.Error {
	currentUserID := ctx.Value(util.ContextKey("user_id"))
	if id != currentUserID && !permission.Has(ctx, permission.PermissionDeleteUser) {
		return errors.Forbidden()
	}

	// Ensure the user exists.
	u, err := uu.repo.Get(id)
	if err != nil {
		return err
	}

	err = uu.repo.Delete(u.GetID())
	if err != nil {
		return err
	}

	return nil
}
