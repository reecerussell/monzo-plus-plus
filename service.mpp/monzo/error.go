package monzo

// Error is used to read monzo error responses.
type Error struct {
	Code             string `json:"code"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
	Message          string `json:"message"`
}
