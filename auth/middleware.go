package auth

import (
	"github.com/yonasadiel/helios"
)

// LoggedInMiddleware check whether user is authenticated or not
// and send errUnauthorized when user is not logged in.
func LoggedInMiddleware(f helios.HTTPHandler) helios.HTTPHandler {
	return func(req helios.Request) {
		var userToken string
		var userSession Session

		userToken, _ = req.GetSessionData(UserTokenSessionKey).(string)

		if userToken == "" {
			req.SendJSON(errUnauthorized.GetMessage(), errUnauthorized.StatusCode)
			return
		}

		helios.DB.Where("token = ?", userToken).Preload("User").First(&userSession)
		if userSession.ID == 0 {
			req.SendJSON(errUnauthorized.GetMessage(), errUnauthorized.StatusCode)
			return
		}

		req.SetContextData(UserContextKey, *userSession.User)
		f(req)
	}
}
