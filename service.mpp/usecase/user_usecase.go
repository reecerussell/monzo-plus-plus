package usecase

import (
	"github.com/reecerussell/monzo-plus-plus/service.mpp/domain/model"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/domain/repository"
	"github.com/reecerussell/monzo-plus-plus/service.mpp/monzo"
)

type UserUsecase interface {
	FindByStateToken(token string) (*model.User, error)
	Get(userID string) (*model.User, error)
	New() (*model.User, error)
	SetAccessToken(u *model.User, ac *monzo.AccessToken) error
	GetAccessToken(userID string) (string, error)
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (uu *userUsecase) FindByStateToken(token string) (*model.User, error) {
	return uu.repo.FindByStateToken(token)
}

func (uu *userUsecase) Get(userID string) (*model.User, error) {
	return uu.repo.Get(userID)
}

func (uu *userUsecase) New() (*model.User, error) {
	u := model.NewUser()

	err := uu.repo.Create(u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

// SetAccessToken updates the user's token values with the new access token.
func (uu *userUsecase) SetAccessToken(u *model.User, ac *monzo.AccessToken) error {
	err := u.UpdateMonzoID(ac.UserID)
	if err != nil {
		return err
	}

	ut := u.GetToken()

	err = ut.UpdateAccessToken(ac.AccessToken)
	if err != nil {
		return err
	}

	err = ut.UpdateRefreshToken(ac.RefreshToken)
	if err != nil {
		return err
	}

	err = ut.UpdateTokenType(ac.TokenType)
	if err != nil {
		return err
	}

	err = ut.UpdateExpires(ac.ExpiresIn)
	if err != nil {
		return err
	}

	err = uu.repo.Update(u)
	if err != nil {
		return err
	}

	return nil
}

// RefreshAccessToken refreshed the given user's access token.
func (uu *userUsecase) RefreshAccessToken(u *model.User) error {
	ac, err := monzo.RefreshAccessToken(u.GetToken().GetRefreshToken())
	if err != nil {
		return err
	}

	return uu.SetAccessToken(u, ac)
}

// GetAccessToken returns a user's access token, but refreshes it if needed.
func (uu *userUsecase) GetAccessToken(userID string) (string, error) {
	u, err := uu.repo.Get(userID)
	if err != nil {
		return "", err
	}

	if !monzo.IsAuthenticated(u) {
		err = uu.RefreshAccessToken(u)
		if err != nil {
			return "", nil
		}
	}

	return u.GetToken().GetAccessToken(), nil
}
