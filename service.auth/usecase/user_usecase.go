package usecase

import (
	"context"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/service"
	"github.com/reecerussell/monzo-plus-plus/service.auth/password"
)

// UserUsecase is a high-level interface providing methods to perform
// CRUD operations, as well as, more specific operations on the User domain.
type UserUsecase interface {
	Create(d *dto.CreateUser) (*dto.User, errors.Error)
	Get(id string) (*dto.User, errors.Error)
	GetList(term string) ([]*dto.User, errors.Error)
	GetPending(term string) ([]*dto.User, errors.Error)
	Update(d *dto.UpdateUser) errors.Error
	Enable(id string) errors.Error
	Delete(id string) errors.Error
	WithUser(ctx context.Context, accessToken string) (context.Context, errors.Error)
}

type userUsecase struct {
	repo repository.UserRepository
	serv *service.UserService
	ps   password.Service
}

func (uu *userUsecase) Create(d *dto.CreateUser) (*dto.User, errors.Error) {
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

func (uu *userUsecase) Get(id string) (*dto.User, errors.Error) {
	u, err := uu.repo.Get(id)
	if err != nil {
		return nil, err
	}

	return u.DTO(), nil
}

func (uu *userUsecase) GetList(term string) ([]*dto.User, errors.Error) {
	users, err := uu.repo.GetList(term)
	if err != nil {
		return nil, err
	}

	return convertToDTOs(users), nil
}

func (uu *userUsecase) GetPending(term string) ([]*dto.User, errors.Error) {
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

func (uu *userUsecase) Update(d *dto.UpdateUser) errors.Error {
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

func (uu *userUsecase) Enable(id string) errors.Error {
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

func (uu *userUsecase) Delete(id string) errors.Error {
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

func (uu *userUsecase) WithUser(ctx context.Context, accessToken string) (context.Context, errors.Error) {
	return ctx, nil
}
