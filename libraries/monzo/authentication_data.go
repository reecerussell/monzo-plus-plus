package monzo

// AuthenticationData is used to read the response body of /ping/whoami requests.
type AuthenticationData struct {
	Authenticated bool   `json:"authenticated"`
	ClientID      string `json:"client_id"`
	UserID        string `json:"user_id"`
}
