package usecase

import (
	"context"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/service"
	"github.com/reecerussell/monzo-plus-plus/service.auth/permission"
)

type RoleUsecase interface {
	Get(ctx context.Context, id string) (*dto.Role, errors.Error)
	GetList(ctx context.Context, term string) ([]*dto.Role, errors.Error)
	Create(ctx context.Context, d *dto.CreateRole) (*dto.Role, errors.Error)
	Update(ctx context.Context, d *dto.Role) errors.Error
	AddPermission(ctx context.Context, d *dto.RolePermission) errors.Error
	RemovePermission(ctx context.Context, d *dto.RolePermission) errors.Error
	Delete(ctx context.Context, id string) errors.Error
}

type roleUsecase struct {
	repo  repository.RoleRepository
	serv  *service.RoleService
	perms repository.PermissionsRepository
}

func NewRoleUsecase(repo repository.RoleRepository, serv *service.RoleService, perms repository.PermissionsRepository) RoleUsecase {
	return &roleUsecase{
		repo:  repo,
		serv:  serv,
		perms: perms,
	}
}

func (ru *roleUsecase) Get(ctx context.Context, id string) (*dto.Role, errors.Error) {
	if !permission.Has(ctx, permission.PermissionGetRole) {
		return nil, errors.Forbidden()
	}

	role, err := ru.repo.Get(id)
	if err != nil {
		return nil, err
	}

	return role.DTO(), nil
}

func (ru *roleUsecase) GetList(ctx context.Context, term string) ([]*dto.Role, errors.Error) {
	if !permission.Has(ctx, permission.PermissionGetRoleList) {
		return nil, errors.Forbidden()
	}

	roles, err := ru.repo.GetList(term)
	if err != nil {
		return nil, err
	}

	dtos := make([]*dto.Role, len(roles))

	for i, r := range roles {
		dtos[i] = r.DTO()
	}

	return dtos, nil
}

func (ru *roleUsecase) Create(ctx context.Context, d *dto.CreateRole) (*dto.Role, errors.Error) {
	if !permission.Has(ctx, permission.PermissionCreateRole) {
		return nil, errors.Forbidden()
	}

	r, err := model.NewRole(d)
	if err != nil {
		return nil, err
	}

	err = ru.serv.ValidateName(r)
	if err != nil {
		return nil, err
	}

	err = ru.repo.Insert(r)
	if err != nil {
		return nil, err
	}

	return r.DTO(), nil
}

func (ru *roleUsecase) Update(ctx context.Context, d *dto.Role) errors.Error {
	if !permission.Has(ctx, permission.PermissionUpdateRole) {
		return errors.Forbidden()
	}

	r, err := ru.repo.Get(d.ID)
	if err != nil {
		return err
	}

	err = r.Update(d)
	if err != nil {
		return err
	}

	err = ru.serv.ValidateName(r)
	if err != nil {
		return err
	}

	err = ru.repo.Update(r)
	if err != nil {
		return err
	}

	return nil
}

func (ru *roleUsecase) AddPermission(ctx context.Context, d *dto.RolePermission) errors.Error {
	if !permission.Has(ctx, permission.PermissionUpdateRole) {
		return errors.Forbidden()
	}

	r, err := ru.repo.Get(d.RoleID)
	if err != nil {
		return err
	}

	p, err := ru.perms.Get(d.PermissionID)
	if err != nil {
		return err
	}

	err = r.AddPermission(p)
	if err != nil {
		return err
	}

	err = ru.repo.Update(r)
	if err != nil {
		return err
	}

	return nil
}

func (ru *roleUsecase) RemovePermission(ctx context.Context, d *dto.RolePermission) errors.Error {
	if !permission.Has(ctx, permission.PermissionUpdateRole) {
		return errors.Forbidden()
	}

	r, err := ru.repo.Get(d.RoleID)
	if err != nil {
		return err
	}

	p, err := ru.perms.Get(d.PermissionID)
	if err != nil {
		return err
	}

	err = r.RemovePermission(p)
	if err != nil {
		return err
	}

	err = ru.repo.Update(r)
	if err != nil {
		return err
	}

	return nil
}

func (ru *roleUsecase) Delete(ctx context.Context, id string) errors.Error {
	if !permission.Has(ctx, permission.PermissionDeleteRole) {
		return errors.Forbidden()
	}

	err := ru.repo.EnsureExists(id)
	if err != nil {
		return err
	}

	err = ru.repo.Delete(id)
	if err != nil {
		return err
	}

	return nil
}
