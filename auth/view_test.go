package auth

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yonasadiel/helios"
)

func TestLoginViewsSuccess(t *testing.T) {
	helios.App.BeforeTest()
	var user User = NewUser("name", "email", "password")
	helios.DB.Create(&user)

	requestData := LoginRequest{Email: "email", Password: "password"}
	req := helios.MockRequest{
		RequestData: requestData,
		SessionData: make(map[string]interface{}),
	}

	LoginView(&req)

	assert.Equal(t, http.StatusOK, req.StatusCode, "Unexpected status code")

	var returnedUser map[string]interface{}
	errUnmarshalling := json.Unmarshal(req.JSONResponse, &returnedUser)
	if errUnmarshalling != nil {
		t.Errorf("Error unmarshalling: %s", errUnmarshalling)
	}

	var userSession Session
	assert.Equal(t, "name", returnedUser["name"], "Wrong Name")
	assert.Equal(t, "email", returnedUser["email"], "Wrong Email")
	userToken, ok := req.GetSessionData(UserTokenSessionKey).(string)
	assert.True(t, ok, "Fail to convert user token")
	assert.NotEmpty(t, userToken, "Session token is empty")
	helios.DB.Where("token = ?", userToken).Find(&userSession)
	assert.Equal(t, user.ID, userSession.UserID, "user ID is not equal")
}

func TestLoginViewWrongUsername(t *testing.T) {
	helios.App.BeforeTest()
	var user User = NewUser("name", "email", "password")
	helios.DB.Create(&user)

	requestData := LoginRequest{Email: "wrong_email", Password: "password"}
	req := helios.MockRequest{
		RequestData: requestData,
		SessionData: make(map[string]interface{}),
	}

	LoginView(&req)

	assert.Equal(t, errWrongUsernamePassword.GetStatusCode(), req.StatusCode, "Unexpected status code")

	var errMessage map[string]interface{}
	errUnmarshalling := json.Unmarshal(req.JSONResponse, &errMessage)
	if errUnmarshalling != nil {
		t.Errorf("Error unmarshalling: %s", errUnmarshalling)
	}

	assert.Equal(t, errWrongUsernamePassword.Code, errMessage["code"], "Wrong Code")
	assert.Equal(t, errWrongUsernamePassword.Message, errMessage["message"], "Wrong Message")
	assert.Empty(t, req.GetSessionData(UserTokenSessionKey), "User is logged in")
}

func TestLoginViewWrongPassword(t *testing.T) {
	helios.App.BeforeTest()
	var user User = NewUser("name", "email", "password")
	helios.DB.Create(&user)

	requestData := LoginRequest{Email: "email", Password: "wrong_password"}
	req := helios.MockRequest{
		RequestData: requestData,
		SessionData: make(map[string]interface{}),
	}

	LoginView(&req)

	assert.Equal(t, errWrongUsernamePassword.GetStatusCode(), req.StatusCode, "Unexpected status code")

	var errMessage map[string]interface{}
	errUnmarshalling := json.Unmarshal(req.JSONResponse, &errMessage)
	if errUnmarshalling != nil {
		t.Errorf("Error unmarshalling: %s", errUnmarshalling)
	}

	assert.Equal(t, errWrongUsernamePassword.Code, errMessage["code"], "Wrong Code")
	assert.Equal(t, errWrongUsernamePassword.Message, errMessage["message"], "Wrong Message")
	assert.Empty(t, req.GetSessionData(UserTokenSessionKey), "User is logged in")
}

func TestLogoutViewsSuccess(t *testing.T) {
	helios.App.BeforeTest()

	token := "random_token"
	user := User{Email: "email"}
	helios.DB.Create(&user)

	session := Session{Token: token, UserID: user.ID}
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
