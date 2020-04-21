package event

// DeleteToken is a domain event used to delete a user's token
// data from the database. This is typically used if the system has
// discovered that the user's token is invalid.
type DeleteToken struct {
	UserID string
}
