package app

// APIError is standardized error of Charon app
type APIError struct {
	StatusCode int
	Code       string
	Message    string
}

// GetMessage Get the message to shown as response body
func (apiError APIError) GetMessage() map[string]string {
	message := make(map[string]string)
	message["code"] = apiError.Code
	message["message"] = apiError.Message

	return message
}
