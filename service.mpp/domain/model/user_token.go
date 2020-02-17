package model

import (
	"fmt"
	"time"
)

type UserToken struct {
	accessToken  string
	refreshToken string
	expires      time.Time
	tokenType    string
}

func (ut *UserToken) GetAccessToken() string {
	return ut.accessToken
}

func (ut *UserToken) GetRefreshToken() string {
	return ut.refreshToken
}

func (ut *UserToken) GetExpiryDate() time.Time {
	return ut.expires
}

func (ut *UserToken) GetTokenType() string {
	return ut.tokenType
}

func (ut *UserToken) UpdateAccessToken(token string) error {
	if token == "" {
		return fmt.Errorf("user token: access token cannot be empty")
	}

	ut.accessToken = token

	return nil
}

func (ut *UserToken) UpdateRefreshToken(token string) error {
	if token == "" {
		return fmt.Errorf("user token: refresh token cannot be empty")
	}

	ut.refreshToken = token

	return nil
}

func (ut *UserToken) UpdateTokenType(tokenType string) error {
	if tokenType == "" {
		return fmt.Errorf("user token: token type cannot be empty")
	}

	ut.tokenType = tokenType

	return nil
}

func (ut *UserToken) UpdateExpires(seconds int) error {
	exp := time.Now()
	exp.Add(time.Second * time.Duration(seconds))

	ut.expires = exp

	return nil
}
