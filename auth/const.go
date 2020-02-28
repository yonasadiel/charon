package auth

import (
	"net/http"

	"github.com/yonasadiel/helios"
)

const (
	hashCost = 10

	// UserTokenSessionKey is the key of session data that store user id
	UserTokenSessionKey = "user"
	userTokenLength     = 16 // length of the token
	userTokenBytes      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	userTokenIdxBits    = 6                       // 6 bits to represent a userToken index
	userTokenIdxMask    = 1<<userTokenIdxBits - 1 // All 1-bits, as many as userTokenIdxBits
	userTokenIdxMax     = 63 / userTokenIdxBits   // # of userToken indices fitting in 63 bits
	// UserContextKey is the key of context data that store user object
	UserContextKey = "user"

	userTypeLocal       = "local"       // the one that organize the exam
	userTypeParticipant = "participant" // the one that taking the exam

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
