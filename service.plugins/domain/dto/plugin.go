package dto

type Plugin struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	DisplayName    string `json:"displayName"`
	Description    string `json:"description"`
	ConsumedBy     int    `json:"consumedBy"`
	ConsumedByUser bool   `json:"consumedByUser"`
}
