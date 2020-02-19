package dto

import "time"

type User struct {
	ID          string     `json:"id"`
	Username    string     `json:"username"`
	DateEnabled *time.Time `json:"dateEnabled,omitempty"`
	Enabled     bool       `json:"enabled"`
}
