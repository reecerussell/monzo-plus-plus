package datamodel

type Plugin struct {
	ID             string
	Name           string
	DisplayName    string
	Description    string
	ConsumedBy     int
	ConsumedByUser bool
}
