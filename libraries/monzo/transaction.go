package monzo

import "time"

// TransactionEventWrapper is used to read webhook transaction event data.
// An example of this is the 'transaction.created' event.
type TransactionEventWrapper struct {
	Type string      `json:"type"`
	Data Transaction `json:"data"`
}

// Transaction contains data or properties related to transaction data.
type Transaction struct {
	AccountID   string     `json:"account_id"`
	Amount      int        `json:"amount"`
	Created     time.Time  `json:"created"`
	Currency    string     `json:"currency"`
	Description string     `json:"description"`
	ID          string     `json:"id"`
	Category    string     `json:"category"`
	IsLoad      bool       `json:"is_load"`
	Settled     *time.Time `json:"settled"`
}
