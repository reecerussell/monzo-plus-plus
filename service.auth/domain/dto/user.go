package dto

import "time"

type User struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	DateEnabled *time.Time `json:"dateEnabled,omitempty"`
	Enabled     bool       `json:"enabled"`
	MonzoLinked bool       `json:"monzoLinked"`
	AccountID   *string    `json:"accountId"`

	Roles []*Role `json:"roles,omitempty"`
}
