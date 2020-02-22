package auth

import "github.com/yonasadiel/helios"

// LoggedInMiddleware check whether user is authenticated or not
// and send errUnauthorized when user is not logged in.
func LoggedInMiddleware(f helios.HTTPHandler) helios.HTTPHandler {
	return func(req helios.Request) {
		var userEmailSessionData interface{} = req.GetSessionData(UserEmailSessionKey)

		if userEmailSessionData == nil {
			req.SendJSON(errUnauthorized.GetMessage(), errUnauthorized.StatusCode)
		} else {
			var user User
			helios.DB.Where("email = ?", userEmailSessionData.(string)).First(&user)
			req.SetContextData(UserContextKey, user)
			f(req)
		}
	}
}
