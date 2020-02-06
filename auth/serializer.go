package auth

// DeserializeLoginRequest deserialize the request data for logging in
func DeserializeLoginRequest(requestData map[string]string) LoginRequest {
	return LoginRequest{
		Email:    requestData["email"],
		Password: requestData["password"],
	}
}
