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
	assert.Equal(t, errUnauthorized.GetStatusCode(), req.StatusCode, "User should be unauthorized")
}

func TestLoggedInMiddlewareUnknownTekon(t *testing.T) {
	helios.App.BeforeTest()
	var blankHandler = func(req helios.Request) {
		req.SendJSON("OK", 200)
	}
	var wrappedHandler helios.HTTPHandler = LoggedInMiddleware(blankHandler)
	var req helios.MockRequest = helios.NewMockRequest()
	req.SetSessionData(UserTokenSessionKey, "unknown_token")
	wrappedHandler(&req)
	assert.Equal(t, errUnauthorized.GetStatusCode(), req.StatusCode, "User should be unauthorized")
}

func TestLoggedInMiddlewareWrongIP(t *testing.T) {
	helios.App.BeforeTest()

	var blankHandler = func(req helios.Request) {
		req.SendJSON("OK", 200)
	}
	var wrappedHandler helios.HTTPHandler = LoggedInMiddleware(blankHandler)

	var user1 User = NewUser("User 1", "user1", "password1")
	var token string = "random_token"

	helios.DB.Create(&user1)
	helios.DB.Create(&Session{Token: token, UserID: user1.ID, IPAddress: "7.1.1.1"})

	var req helios.MockRequest = helios.NewMockRequest()
	req.SetSessionData(UserTokenSessionKey, token)
	req.RemoteAddr = "7.1.1.2"
	wrappedHandler(&req)

	assert.Equal(t, errUnauthorized.GetStatusCode(), req.StatusCode, "User should be unauthorized")
}

func TestLoggedInMiddlewareAuthorized(t *testing.T) {
	helios.App.BeforeTest()
	var blankHandler = func(req helios.Request) {
		req.SendJSON("OK", 200)
	}
	var wrappedHandler helios.HTTPHandler = LoggedInMiddleware(blankHandler)

	var user1 User = NewUser("User 1", "user1", "password1")
	var token string = "random_token"

	helios.DB.Create(&user1)
	helios.DB.Create(&Session{Token: token, UserID: user1.ID, IPAddress: "7.1.1.1"})

	var req helios.MockRequest = helios.NewMockRequest()
	req.SetSessionData(UserTokenSessionKey, token)
	req.RemoteAddr = "7.1.1.1"
	wrappedHandler(&req)

	assert.Equal(t, 200, req.StatusCode, "User should be authorized")
	userReturned, successCoversion := req.GetContextData(UserTokenSessionKey).(User)
	assert.True(t, successCoversion, "Failed to convert user in context data to user object")
	assert.Equal(t, user1.ID, userReturned.ID, "User object should be on the context data")
}
