package event

// DisablePluginForUser is an event used to disable a plugin for a user.
type DisablePluginForUser struct {
	PluginID, UserID string
}
