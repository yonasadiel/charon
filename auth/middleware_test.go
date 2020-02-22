package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yonasadiel/helios"
)

func TestLoggedInMiddlewareUnauthorized(t *testing.T) {
	helios.App.BeforeTest()
	var blankHandler = func(req helios.Request) {
		req.SendJSON("OK", 200)
	}
	var wrappedHandler helios.HTTPHandler = LoggedInMiddleware(blankHandler)

	var req helios.MockRequest = helios.NewMockRequest()
	wrappedHandler(&req)
	assert.Equal(t, errUnauthorized.StatusCode, req.StatusCode, "User should be unauthorized")
}

func TestLoggedInMiddlewareAuthorized(t *testing.T) {
	helios.App.BeforeTest()
	var blankHandler = func(req helios.Request) {
		req.SendJSON("OK", 200)
	}
	var wrappedHandler helios.HTTPHandler = LoggedInMiddleware(blankHandler)

	var user1 User = NewUser("User 1", "user1", "password1")
	var req helios.MockRequest = helios.NewMockRequest()
	helios.DB.Create(&user1)
	req.SetSessionData(UserEmailSessionKey, user1.Email)
	wrappedHandler(&req)
	assert.Equal(t, 200, req.StatusCode, "User should be authorized")
}
