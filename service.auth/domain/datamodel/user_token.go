package datamodel

import "time"

type UserToken struct {
	AccessToken  string
	RefreshToken string
	Expires      time.Time
	TokenType    string
}
