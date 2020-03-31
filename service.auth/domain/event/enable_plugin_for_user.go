package event

// EnablePluginForUser is an event used to enable a plugin for a user.
type EnablePluginForUser struct {
	PluginID, UserID string
}
