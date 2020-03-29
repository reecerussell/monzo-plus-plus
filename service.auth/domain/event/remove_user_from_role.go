package event

// RemoveUserFromRole is an event used to unassign a user from a role. The
// event object is used to store the user's and role's ids.
type RemoveUserFromRole struct {
	UserID, RoleID string
}
