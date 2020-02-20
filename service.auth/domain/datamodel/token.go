package datamodel

import "time"

type Token struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	Expires      time.Time
	TokenType    string
}
