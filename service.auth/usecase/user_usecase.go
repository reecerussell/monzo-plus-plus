package usecase

import (
	"context"
	"strings"

	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"

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
	AddToRole(ctx context.Context, d *dto.UserRole) errors.Error
	RemoveFromRole(ctx context.Context, d *dto.UserRole) errors.Error
	GetRoles(ctx context.Context, id string) ([]*dto.Role, errors.Error)
	GetAvailableRoles(ctx context.Context, id string) ([]*dto.Role, errors.Error)
	EnablePlugin(ctx context.Context, d *dto.UserPlugin) errors.Error
	DisablePlugin(ctx context.Context, d *dto.UserPlugin) errors.Error
	GetAccounts(ctx context.Context, id string) ([]*monzo.AccountData, errors.Error)
	SetAccount(ctx context.Context, d *dto.UserAccount) errors.Error
	Delete(ctx context.Context, id string) errors.Error
}

type userUsecase struct {
	repo  repository.UserRepository
	roles repository.RoleRepository
	serv  *service.UserService
	ps    password.Service
	auth  UserAuthUsecase
}

// NewUserUsecase instantiates a new instance of UserUsecase with the given dependencies.
func NewUserUsecase(repo repository.UserRepository,
	roles repository.RoleRepository,
	serv *service.UserService,
	ps password.Service,
	auth UserAuthUsecase) UserUsecase {
	return &userUsecase{
		repo:  repo,
		roles: roles,
		serv:  serv,
		ps:    ps,
		auth:  auth,
	}
}

func (uu *userUsecase) Create(ctx context.Context, d *dto.CreateUser) (*dto.User, errors.Error) {
	if !permission.Has(ctx, permission.PermissionCreateUser) {
		return nil, errors.Forbidden()
	}

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

func (uu *userUsecase) AddToRole(ctx context.Context, d *dto.UserRole) errors.Error {
	if !permission.Has(ctx, permission.PermissionRoleManager) {
		return errors.Forbidden()
	}

	uc, rc := make(chan *model.User, 1), make(chan *model.Role, 1)

	var eg errors.Group
	eg.Go(func() errors.Error {
		u, err := uu.repo.Get(d.UserID)
		if err != nil {
			return err
		}

		uc <- u

		return nil
	})
	eg.Go(func() errors.Error {
		r, err := uu.roles.Get(d.RoleID)
		if err != nil {
			return err
		}

		rc <- r

		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	u, r := <-uc, <-rc

	err := u.AddToRole(r)
	if err != nil {
		return err
	}

	err = uu.repo.Update(u)
	if err != nil {
		return err
	}

	return nil
}

func (uu *userUsecase) RemoveFromRole(ctx context.Context, d *dto.UserRole) errors.Error {
	if !permission.Has(ctx, permission.PermissionRoleManager) {
		return errors.Forbidden()
	}

	uc, rc := make(chan *model.User, 1), make(chan *model.Role, 1)

	var eg errors.Group
	eg.Go(func() errors.Error {
		u, err := uu.repo.Get(d.UserID)
		if err != nil {
			return err
		}

		uc <- u

		return nil
	})
	eg.Go(func() errors.Error {
		r, err := uu.roles.Get(d.RoleID)
		if err != nil {
			return err
		}

		rc <- r

		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	u, r := <-uc, <-rc

	err := u.RemoveFromRole(r)
	if err != nil {
		return err
	}

	err = uu.repo.Update(u)
	if err != nil {
		return err
	}

	return nil
}

func (uu *userUsecase) GetRoles(ctx context.Context, id string) ([]*dto.Role, errors.Error) {
	if ctx.Value(util.ContextKey("user_id")) != id &&
		!permission.Has(ctx, permission.PermissionRoleManager) {
		return nil, errors.Forbidden()
	}

	roles, err := uu.roles.GetForUser(id)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.Role, len(roles))

	for i, r := range roles {
		dtos[i] = r.DTO()
	}

	return dtos, nil
}

func (uu *userUsecase) GetAvailableRoles(ctx context.Context, id string) ([]*dto.Role, errors.Error) {
	if !permission.Has(ctx, permission.PermissionRoleManager) {
		return nil, errors.Forbidden()
	}

	roles, err := uu.roles.GetAvailableForUser(id)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.Role, len(roles))

	for i, r := range roles {
		dtos[i] = r.DTO()
	}

	return dtos, nil
}

// EnablePlugin enables a plugin for a specific user. This method can only be executed
// if the current user is the target user or has the PluginManager permission.
func (uu *userUsecase) EnablePlugin(ctx context.Context, d *dto.UserPlugin) errors.Error {
	if ctx.Value(util.ContextKey("user_id")) != d.UserID &&
		!permission.Has(ctx, permission.PermissionPluginManager) {
		return errors.Forbidden()
	}

	u, err := uu.repo.Get(d.UserID)
	if err != nil {
		return err
	}

	u.EnablePlugin(d.PluginID)

	err = uu.repo.Update(u)
	if err != nil {
		return err
	}

	return nil
}

// DisablePlugin disables a plugin for a specific user. This method can only be executed
// if the current user is the target user or has the PluginManager permission.
func (uu *userUsecase) DisablePlugin(ctx context.Context, d *dto.UserPlugin) errors.Error {
	if ctx.Value(util.ContextKey("user_id")) != d.UserID &&
		!permission.Has(ctx, permission.PermissionPluginManager) {
		return errors.Forbidden()
	}

	u, err := uu.repo.Get(d.UserID)
	if err != nil {
		return err
	}

	u.DisablePlugin(d.PluginID)

	err = uu.repo.Update(u)
	if err != nil {
		return err
	}

	return nil
}

func (uu *userUsecase) GetAccounts(ctx context.Context, id string) ([]*monzo.AccountData, errors.Error) {
	currentUserID := ctx.Value(util.ContextKey("user_id"))
	if id != currentUserID && !permission.Has(ctx, permission.PermissionUpdateUser) {
		return nil, errors.Forbidden()
	}

	ac, err := uu.auth.GetMonzoAccessToken(id)
	if err != nil {
		return nil, err
	}

	accounts, aErr := monzo.Accounts(ac)
	if aErr != nil {
		if strings.Contains(aErr.Error(), monzo.ErrForbidden.Text()) {
			return nil, errors.BadRequest("Open your Monzo app to allow Monzo++ to get your accounts.")
		}

		return nil, errors.InternalError(aErr)
	}

	return accounts, nil
}

// SetAccount is used to configure a user's desired account. This will
// also register a webhook with the given account.
func (uu *userUsecase) SetAccount(ctx context.Context, d *dto.UserAccount) errors.Error {
	currentUserID := ctx.Value(util.ContextKey("user_id"))
	if d.UserID != currentUserID && !permission.Has(ctx, permission.PermissionUpdateUser) {
		return errors.Forbidden()
	}

	ac, err := uu.auth.GetMonzoAccessToken(d.UserID)
	if err != nil {
		return err
	}

	u, err := uu.repo.Get(d.UserID)
	if err != nil {
		return err
	}

	err = u.UpdateAccountID(d.AccountID, ac)
	if err != nil {
		if strings.Contains(err.Text(), monzo.ErrForbidden.Text()) {
			return errors.BadRequest("Open your Monzo app to allow Monzo++ to get your accounts.")
		}

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
