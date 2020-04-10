package auth

// LoginRequest is request for logging in
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserResponse is JSON representation of User.
type UserResponse struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// SerializeUser serialize user to UserResponse
func SerializeUser(user User) UserResponse {
	return UserResponse{
		Username: user.Username,
		Name:     user.Name,
		Role:     user.Role,
	}
}
