package auth

import (
	"net/http"

	"github.com/yonasadiel/helios"
)

// LoginView logging in the user
func LoginView(req helios.Request) {
	var loginRequest LoginRequest
	req.DeserializeRequestData(&loginRequest)
	userSession, err := Login(loginRequest.Username, loginRequest.Password, req.ClientIP())

	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
	} else {
		req.SetSessionData(UserTokenSessionKey, userSession.Token)
		req.SaveSession()
		req.SendJSON(SerializeUser(*userSession.User), http.StatusOK)
	}
}

// LogoutView clear user session data
func LogoutView(req helios.Request) {
	var user User = req.GetContextData(UserTokenSessionKey).(User)
	Logout(user)
	req.SetSessionData(UserTokenSessionKey, "")
	req.SendJSON(nil, http.StatusOK)
}

// UserListView returns all the users
func UserListView(req helios.Request) {
	var user User = req.GetContextData(UserTokenSessionKey).(User)
	var users []User = GetAllUser(user)

	serializedUsers := make([]UserData, 0)
	for _, user := range users {
		serializedUsers = append(serializedUsers, SerializeUser(user))
	}
	req.SendJSON(serializedUsers, http.StatusOK)
}

// UserCreateView creates the user
func UserCreateView(req helios.Request) {
	user, ok := req.GetContextData(UserContextKey).(User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var userData UserWithPasswordData
	var newUser User
	var err helios.Error
	err = req.DeserializeRequestData(&userData)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	err = DeserializeUserWithPassword(userData, &newUser)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	newUser.ID = 0
	UpsertUser(user, &newUser)
	req.SendJSON(SerializeUser(newUser), http.StatusCreated)
}

// UserSessionLockView locks the session of user
func UserSessionLockView(req helios.Request) {
	user, ok := req.GetContextData(UserContextKey).(User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var username string = req.GetURLParam("username")
	var err helios.Error = LockUserSession(user, username)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	req.SendJSON("OK", http.StatusOK)
}

// UserSessionUnlockView locks the session of user
func UserSessionUnlockView(req helios.Request) {
	user, ok := req.GetContextData(UserContextKey).(User)
	if !ok {
		req.SendJSON(helios.ErrInternalServerError.GetMessage(), helios.ErrInternalServerError.GetStatusCode())
		return
	}

	var username string = req.GetURLParam("username")
	var err helios.Error = UnlockUserSession(user, username)
	if err != nil {
		req.SendJSON(err.GetMessage(), err.GetStatusCode())
		return
	}
	req.SendJSON("OK", http.StatusOK)
}
