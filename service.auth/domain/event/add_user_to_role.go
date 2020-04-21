package event

// AddUserToRole is an event used to assign a user to a role. The
// event object is used to store the user's and role's ids.
type AddUserToRole struct {
	UserID, RoleID string
}
