package dto

// ChangePassword is used to read a request body, which is then
// used to change a user's password.
type ChangePassword struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}
