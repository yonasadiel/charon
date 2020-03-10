package auth

// LoginRequest is request for logging in
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserResponse is JSON representation of User.
type UserResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

// SerializeUser serialize user to UserResponse
func SerializeUser(user User) UserResponse {
	return UserResponse{
		Email: user.Email,
		Name:  user.Name,
		Role:  user.Role,
	}
}
