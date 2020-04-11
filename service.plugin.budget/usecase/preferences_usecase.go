package usecase

import (
	"context"

	"github.com/reecerussell/monzo-plus-plus/libraries/errors"
	"github.com/reecerussell/monzo-plus-plus/libraries/permission"
	"github.com/reecerussell/monzo-plus-plus/libraries/util"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/domain/repository"
)

type PreferencesUsecase interface {
	Get(ctx context.Context, userID string) (*dto.Preferences, errors.Error)
	GetMonthlyBudget(userID string) (int, errors.Error)
	Update(ctx context.Context, d *dto.Preferences) errors.Error
}

type preferencesUsecase struct {
	repo repository.PreferencesRepository
}

func NewPreferencesUsecase(repo repository.PreferencesRepository) PreferencesUsecase {
	return &preferencesUsecase{
		repo: repo,
	}
}

func (pu *preferencesUsecase) Get(ctx context.Context, userID string) (*dto.Preferences, errors.Error) {
	cui := ctx.Value(util.ContextKey("user_id")).(string)
	if cui != userID && !permission.Has(ctx, permission.PermissionGetUser) {
		return nil, errors.Forbidden()
	}

	p, err := pu.repo.Get(userID)
	if err != nil {
		return nil, errors.InternalError(err)
	}

	return &dto.Preferences{
		UserID:        p.GetUserID(),
		MonthlyBudget: p.GetMonthlyBudget(),
	}, nil
}

func (pu *preferencesUsecase) GetMonthlyBudget(userID string) (int, errors.Error) {
	p, err := pu.repo.Get(userID)
	if err != nil {
		return 0, errors.InternalError(err)
	}

	return p.GetMonthlyBudget(), nil
}

func (pu *preferencesUsecase) Update(ctx context.Context, d *dto.Preferences) errors.Error {
	cui := ctx.Value(util.ContextKey("user_id")).(string)
	if cui != d.UserID && !permission.Has(ctx, permission.PermissionUpdateUser) {
		return errors.Forbidden()
	}

	p := model.NewPreferences(d.UserID)
	err := p.UpdateMonthlyBudget(d.MonthlyBudget)
	if err != nil {
		return err
	}

	if err := pu.repo.Update(p); err != nil {
		return errors.InternalError(err)
	}

	return nil
}
