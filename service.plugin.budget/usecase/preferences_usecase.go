package usecase

import (
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/domain/dto"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.plugin.budget/domain/repository"
)

type PreferencesUsecase interface {
	Get(userID string) (*dto.Preferences, error)
	GetMonthlyBudget(userID string) (int, error)
	Update(d *dto.Preferences) error
}

type preferencesUsecase struct {
	repo repository.PreferencesRepository
}

func NewPreferencesUsecase(repo repository.PreferencesRepository) PreferencesUsecase {
	return &preferencesUsecase{
		repo: repo,
	}
}

func (pu *preferencesUsecase) Get(userID string) (*dto.Preferences, error) {
	p, err := pu.repo.Get(userID)
	if err != nil {
		return nil, err
	}

	return &dto.Preferences{
		UserID:        p.GetUserID(),
		MonthlyBudget: p.GetMonthlyBudget(),
	}, nil
}

func (pu *preferencesUsecase) GetMonthlyBudget(userID string) (int, error) {
	p, err := pu.Get(userID)
	if err != nil {
		return 0, err
	}

	return p.MonthlyBudget, nil
}

func (pu *preferencesUsecase) Update(d *dto.Preferences) error {
	p := model.NewPreferences(d.UserID)
	err := p.UpdateMonthlyBudget(d.MonthlyBudget)
	if err != nil {
		return err
	}

	return pu.repo.Update(p)
}
