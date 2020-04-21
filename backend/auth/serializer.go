package auth

import "github.com/yonasadiel/helios"

// LoginRequest is request for logging in
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UserData is JSON representation of User.
type UserData struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

// UserWithPasswordData is JSON representation of User according to auth.UserData,
// including password. Used for synchronization.
type UserWithPasswordData struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Password string `json:"password"`
}

// SerializeUser serialize user to UserData
func SerializeUser(user User) UserData {
	var role string
	if user.IsAdmin() {
		role = "admin"
	} else if user.IsOrganizer() {
		role = "organizer"
	} else if user.IsLocal() {
		role = "local"
	} else if user.IsParticipant() {
		role = "participant"
	}
	return UserData{
		Username: user.Username,
		Name:     user.Name,
		Role:     role,
	}
}

// DeserializeUser deserialize UserData to User
func DeserializeUser(userData UserData, user *User) helios.Error {
	var err helios.ErrorForm = helios.NewErrorForm()
	user.Name = userData.Name
	user.Username = userData.Username
	if user.Name == "" {
		err.FieldError["name"] = helios.ErrorFormFieldAtomic{"Name can't be empty"}
	}
	if user.Username == "" {
		err.FieldError["username"] = helios.ErrorFormFieldAtomic{"Username can't be empty"}
	}
	if userData.Role == "admin" {
		user.Role = UserRoleAdmin
	} else if userData.Role == "organizer" {
		user.Role = UserRoleOrganizer
	} else if userData.Role == "local" {
		user.Role = UserRoleLocal
	} else if userData.Role == "participant" {
		user.Role = UserRoleParticipant
	} else if userData.Role == "" {
		err.FieldError["role"] = helios.ErrorFormFieldAtomic{"Role can't be empty"}
	} else {
		err.FieldError["role"] = helios.ErrorFormFieldAtomic{"Role should be either admin, organizer, local, or participant"}
	}
	if err.IsError() {
		return err
	}
	return nil
}

// SerializeUserWithPassword serialize user to UserWithPasswordData
func SerializeUserWithPassword(user User) UserWithPasswordData {
	var userData UserData = SerializeUser(user)
	return UserWithPasswordData{
		Username: userData.Username,
		Name:     userData.Name,
		Role:     userData.Role,
		Password: user.Password,
	}
}

// DeserializeUserWithPassword deserialize UserWithPasswordData to User
func DeserializeUserWithPassword(userData UserWithPasswordData, user *User) helios.Error {
	err := DeserializeUser(UserData{
		Name:     userData.Name,
		Username: userData.Username,
		Role:     userData.Role,
	}, user)
	user.Password = userData.Password
	if err != nil {
		return err
	}
	return nil
}
