package auth

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yonasadiel/helios"
)

func TestLoginView(t *testing.T) {
	helios.App.BeforeTest()

	type loginViewTestCase struct {
		user               User
		username           string
		password           string
		expectedStatusCode int
		expectedError      helios.Error
	}
	testCases := []loginViewTestCase{{
		user:               UserFactorySaved(User{Username: "user1", Password: "password"}),
		username:           "user1",
		password:           "password",
		expectedStatusCode: http.StatusOK,
	}, {
		user:               UserFactorySaved(User{Username: "user2", Password: "password"}),
		username:           "wrong_username",
		password:           "password",
		expectedStatusCode: errWrongUsernamePassword.StatusCode,
		expectedError:      errWrongUsernamePassword,
	}, {
		user:               UserFactorySaved(User{Username: "user3", Password: "password"}),
		username:           "user3",
		password:           "wrong_password",
		expectedStatusCode: errWrongUsernamePassword.StatusCode,
		expectedError:      errWrongUsernamePassword,
	}}
	for i, testCase := range testCases {
		t.Logf("Test LoginView testcase: %d", i)
		requestData := LoginRequest{Username: testCase.username, Password: testCase.password}
		req := helios.MockRequest{
			RequestData: requestData,
			SessionData: make(map[string]interface{}),
		}
		LoginView(&req)
		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
		if testCase.expectedError == nil {
			var errMarshall error
			var userSerialized []byte
			var ok bool
			var userToken string
			var userSession Session
			userSerialized, errMarshall = json.Marshal(SerializeUser(testCase.user))
			userToken, ok = req.GetSessionData(UserTokenSessionKey).(string)
			helios.DB.Where("token = ?", userToken).Find(&userSession)
			assert.Nil(t, errMarshall)
			assert.True(t, ok)
			assert.Equal(t, userSerialized, req.JSONResponse)
			assert.Equal(t, testCase.user.ID, userSession.UserID)
			assert.NotEmpty(t, userToken)
		} else {
			var errMarshall error
			var errSerialized []byte
			errSerialized, errMarshall = json.Marshal(testCase.expectedError.GetMessage())
			assert.Nil(t, errMarshall)
			assert.Empty(t, req.GetSessionData(UserTokenSessionKey), "Session data should remain empty")
			assert.Equal(t, errSerialized, req.JSONResponse)
		}
	}
}

func TestLogoutView(t *testing.T) {
	helios.App.BeforeTest()

	var user User = UserFactorySaved(User{Username: "username", Password: "password"})
	var token string = "random_token"
	var session Session = Session{Token: token, UserID: user.ID}

	helios.DB.Create(&session)

	sessionData := make(map[string]interface{})
	sessionData[UserTokenSessionKey] = token
	contextData := make(map[string]interface{})
	contextData[UserContextKey] = user

	req := helios.MockRequest{
		RequestData: nil,
		SessionData: sessionData,
		ContextData: contextData,
	}

	LogoutView(&req)

	assert.Equal(t, http.StatusOK, req.StatusCode, "Unexpected status code")
	assert.Empty(t, sessionData[UserTokenSessionKey], "User token should be removed")

	var userSession Session
	helios.DB.Where("token = ?", token).First(&userSession)
	assert.Equal(t, uint(0), userSession.ID, "User session should be deleted")
}
