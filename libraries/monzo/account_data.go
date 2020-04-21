package monzo

import "time"

type AccountData struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	Created     time.Time `json:"created"`
}

type AccountList struct {
	Accounts []*AccountData `json:"accounts"`
}
