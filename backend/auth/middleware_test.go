package auth

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yonasadiel/helios"
)

func TestLoggedInMiddleware(t *testing.T) {
	helios.App.BeforeTest()

	var user User = UserFactorySaved(User{})
	var token string = "auth_token"
	var blankHandler = func(req helios.Request) {
		req.SendJSON("OK", http.StatusOK)
	}
	var wrappedHandler helios.HTTPHandler = LoggedInMiddleware(blankHandler)
	helios.DB.Create(&Session{Token: token, UserID: user.ID, IPAddress: "7.1.1.1"})
	type loggedInMiddlewareTestCase struct {
		sessionToken       string
		expectedStatusCode int
		remoteAddr         string
	}
	testCases := []loggedInMiddlewareTestCase{{
		expectedStatusCode: errUnauthorized.StatusCode,
	}, {
		sessionToken:       "unknown_token",
		expectedStatusCode: errUnauthorized.StatusCode,
	}, {
		sessionToken:       token,
		expectedStatusCode: errUnauthorized.StatusCode,
		remoteAddr:         "7.1.1.2",
	}, {
		sessionToken:       token,
		expectedStatusCode: http.StatusOK,
		remoteAddr:         "7.1.1.1",
	}}
	for i, testCase := range testCases {
		t.Logf("Test LoggedInMiddleware testcase: %d", i)
		var req helios.MockRequest

		req = helios.NewMockRequest()
		req.SetSessionData(UserTokenSessionKey, testCase.sessionToken)
		req.RemoteAddr = testCase.remoteAddr
		wrappedHandler(&req)

		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
		if testCase.expectedStatusCode == http.StatusOK {
			userReturned, successCoversion := req.GetContextData(UserTokenSessionKey).(User)
			assert.True(t, successCoversion, "Failed to convert user in context data to user object")
			assert.Equal(t, user.ID, userReturned.ID, "User object should be on the context data")
		}
	}
}
