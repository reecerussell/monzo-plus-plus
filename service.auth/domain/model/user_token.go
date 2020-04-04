package model

import (
	"time"

	"github.com/reecerussell/monzo-plus-plus/libraries/monzo"
	"github.com/reecerussell/monzo-plus-plus/service.auth/domain/datamodel"
)

type UserToken struct {
	accessToken  string
	refreshToken string
	expires      time.Time
	tokenType    string
}

func NewUserToken(d *monzo.AccessToken) *UserToken {
	expiry := time.Now().UTC().Add(time.Second * time.Duration(d.ExpiresIn))

	return &UserToken{
		accessToken:  d.AccessToken,
		refreshToken: d.RefreshToken,
		expires:      expiry,
		tokenType:    d.TokenType,
	}
}

func (ut *UserToken) DataModel() *datamodel.UserToken {
	return &datamodel.UserToken{
		AccessToken:  ut.accessToken,
		RefreshToken: ut.refreshToken,
		Expires:      ut.expires,
		TokenType:    ut.tokenType,
	}
}

func UserTokenFromDataModal(d *datamodel.UserToken) *UserToken {
	return &UserToken{
		accessToken:  d.AccessToken,
		refreshToken: d.RefreshToken,
		expires:      d.Expires,
		tokenType:    d.TokenType,
	}
}
