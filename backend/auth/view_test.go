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

	assert.Equal(t, http.StatusOK, req.StatusCode)
	assert.Empty(t, sessionData[UserTokenSessionKey], "User token should be removed")

	var userSession Session
	helios.DB.Where("token = ?", token).First(&userSession)
	assert.Equal(t, uint(0), userSession.ID, "User session should be deleted")
}

func TestUserListView(t *testing.T) {
	helios.App.BeforeTest()

	var user User = UserFactorySaved(User{})
	var req helios.MockRequest
	req = helios.NewMockRequest()
	req.SetContextData(UserContextKey, user)
	UserListView(&req)

	assert.Equal(t, http.StatusOK, req.StatusCode)
}

func TestUserCreateView(t *testing.T) {
	helios.App.BeforeTest()

	var userLocal User = UserFactorySaved(User{Role: UserRoleLocal})
	type userCreateTestCase struct {
		user               interface{}
		requestData        string
		expectedStatusCode int
		expectedErrorCode  string
		expectedUserCount  int
	}
	testCases := []userCreateTestCase{{
		user:               userLocal,
		requestData:        `{"id":2,"name":"User 1","username":"user1","role":"participant"}`,
		expectedStatusCode: http.StatusCreated,
		expectedUserCount:  2,
	}, {
		user:               userLocal,
		requestData:        `{"name":""}`,
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  "form_error",
		expectedUserCount:  2,
	}, {
		user:               "bad_user",
		requestData:        `{"name":""}`,
		expectedStatusCode: http.StatusInternalServerError,
		expectedErrorCode:  helios.ErrInternalServerError.Code,
		expectedUserCount:  2,
	}, {
		user:               userLocal,
		requestData:        `bad_request_data`,
		expectedStatusCode: http.StatusBadRequest,
		expectedErrorCode:  helios.ErrJSONParseFailed.Code,
		expectedUserCount:  2,
	}}
	for i, testCase := range testCases {
		t.Logf("Test UserCreateView testcase: %d", i)
		var eventCount int
		var req helios.MockRequest
		req = helios.NewMockRequest()
		req.SetContextData(UserContextKey, testCase.user)
		req.RequestData = testCase.requestData

		UserCreateView(&req)

		helios.DB.Model(User{}).Count(&eventCount)
		assert.Equal(t, testCase.expectedStatusCode, req.StatusCode)
		assert.Equal(t, testCase.expectedUserCount, eventCount)
		if testCase.expectedErrorCode != "" {
			var err map[string]interface{}
			json.Unmarshal(req.JSONResponse, &err)
			assert.Equal(t, testCase.expectedErrorCode, err["code"])
		}
	}
}
