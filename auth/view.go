package auth

import (
	"net/http"

	"github.com/yonasadiel/helios"
)

// LoginView logging in the user
func LoginView(req helios.Request) {
	var loginRequest LoginRequest
	req.DeserializeRequestData(&loginRequest)
	userSession, err := Login(loginRequest.Email, loginRequest.Password, req.ClientIP())

	if err != nil {
		req.SendJSON(err.GetMessage(), err.StatusCode)
	} else {
		req.SetSessionData(UserTokenSessionKey, userSession.Token)
		req.SendJSON(SerializeUser(*userSession.User), http.StatusOK)
	}
}

// LogoutView clear user session data
func LogoutView(req helios.Request) {
	var user User

	user = req.GetContextData(UserTokenSessionKey).(User)
	Logout(user)
	req.SetSessionData(UserTokenSessionKey, "")
	req.SendJSON(nil, http.StatusOK)
}
