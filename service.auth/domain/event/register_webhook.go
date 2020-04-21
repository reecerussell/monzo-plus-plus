package event

// RegisterWebhook is an event used to register a webhook with a
// user's Monzo account once they've completed their registration.
type RegisterWebhook struct {
	UserID      string
	AccessToken string
	AccountID   string
}
