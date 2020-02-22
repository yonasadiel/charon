package auth

import (
	"net/http"

	"github.com/yonasadiel/helios"
)

const (
	hashCost = 10

	// UserEmailSessionKey is the key of session data that store user id
	UserEmailSessionKey = "user"
	// UserContextKey is the key of context data that store user object
	UserContextKey = "user"
)

/*** APIError of auth package ***/
var errWrongUsernamePassword = helios.APIError{
	StatusCode: http.StatusBadRequest,
	Code:       "login_wrong_email_or_password",
	Message:    "Wrong email / password",
}

var errUnauthorized = helios.APIError{
	StatusCode: http.StatusUnauthorized,
	Code:       "unauthorized",
	Message:    "You need to log in first",
}
