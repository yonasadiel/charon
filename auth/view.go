package auth

import (
	"net/http"

	"github.com/yonasadiel/helios"
)

// LoginView logging in the user
func LoginView(req helios.Request) {
	var loginRequest LoginRequest
	req.DeserializeRequestData(&loginRequest)
	user, err := Login(loginRequest)

	if err != nil {
		req.SendJSON(err.GetMessage(), err.StatusCode)
	} else {
		req.SetSessionData(UserEmailSessionKey, user.Email)
		req.SendJSON(SerializeUser(*user), http.StatusOK)
	}
}

// LogoutView clear user session data
func LogoutView(req helios.Request) {
	req.SetSessionData(UserEmailSessionKey, "")
	req.SendJSON(nil, http.StatusOK)
}
