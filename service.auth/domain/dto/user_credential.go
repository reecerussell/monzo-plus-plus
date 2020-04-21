package dto

// UserCredential contains a user's username and password,
// which can be used to verify a user in the OAuth flow.
type UserCredential struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
