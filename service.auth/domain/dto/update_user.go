package dto

// UpdateUser is a data-transfer object that provides the user domain
// model with only, and exactly what it needs to update a user.
type UpdateUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
