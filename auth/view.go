package auth

import (
	"net/http"

	"github.com/yonasadiel/charon/app"
)

// LoginView logging in the user
func LoginView(req app.Request) {
	var requestData map[string]string = req.GetRequestData()
	var loginRequest LoginRequest = DeserializeLoginRequest(requestData)
	user, err := Login(loginRequest)

	if err != nil {
		req.SendJSON(err.GetMessage(), err.StatusCode)
	} else {
		req.SetSessionData(UserIDSessionKey, user)
		req.SendJSON(user, http.StatusOK)
	}
}
