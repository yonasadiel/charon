package auth

import (
	"net/http"

	"github.com/yonasadiel/helios"
)

const (
	hashCost = 10

	// UserTokenSessionKey is the key of session data that store user id
	UserTokenSessionKey = "user"
	// UserContextKey is the key of context data that store user object
	UserContextKey = "user"

	// UserRoleAdmin is the administrator of the website
	UserRoleAdmin = 40
	// UserRoleOrganizer is the one that organize all the locals
	UserRoleOrganizer = 30
	// UserRoleLocal is the one that organize the local exam
	UserRoleLocal = 20
	// UserRoleParticipant is the one that taking the exam
	UserRoleParticipant = 10

	userTokenLength  = 16 // length of the token
	userTokenBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	userTokenIdxBits = 6                       // 6 bits to represent a userToken index
	userTokenIdxMask = 1<<userTokenIdxBits - 1 // All 1-bits, as many as userTokenIdxBits
	userTokenIdxMax  = 63 / userTokenIdxBits   // # of userToken indices fitting in 63 bits
)

var errWrongUsernamePassword = helios.ErrorForm{
	Code:          "login_wrong_username_or_password",
	NonFieldError: helios.ErrorFormFieldAtomic{"Wrong username / password"},
}

var errUnauthorized = helios.ErrorAPI{
	StatusCode: http.StatusUnauthorized,
	Code:       "unauthorized",
	Message:    "You need to log in first",
}

var errUserRoleTooHigh = helios.ErrorAPI{
	StatusCode: http.StatusBadRequest,
	Code:       "user_role_too_high",
	Message:    "Can't create user with same or higher role",
}
