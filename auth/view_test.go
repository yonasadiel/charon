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

	assert.Equal(t, "name", returnedUser["name"], "Wrong Name")
	assert.Equal(t, "email", returnedUser["email"], "Wrong Email")
	assert.Equal(t, req.GetSessionData(UserEmailSessionKey), user.Email, "Session is not changed")
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

	assert.Equal(t, errWrongUsernamePassword.StatusCode, req.StatusCode, "Unexpected status code")

	var errMessage map[string]interface{}
	errUnmarshalling := json.Unmarshal(req.JSONResponse, &errMessage)
	if errUnmarshalling != nil {
		t.Errorf("Error unmarshalling: %s", errUnmarshalling)
	}

	assert.Equal(t, errWrongUsernamePassword.Code, errMessage["code"], "Wrong Code")
	assert.Equal(t, errWrongUsernamePassword.Message, errMessage["message"], "Wrong Message")
	assert.Equal(t, req.GetSessionData(UserEmailSessionKey), nil, "User is logged in")
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

	assert.Equal(t, errWrongUsernamePassword.StatusCode, req.StatusCode, "Unexpected status code")

	var errMessage map[string]interface{}
	errUnmarshalling := json.Unmarshal(req.JSONResponse, &errMessage)
	if errUnmarshalling != nil {
		t.Errorf("Error unmarshalling: %s", errUnmarshalling)
	}

	assert.Equal(t, errWrongUsernamePassword.Code, errMessage["code"], "Wrong Code")
	assert.Equal(t, errWrongUsernamePassword.Message, errMessage["message"], "Wrong Message")
	assert.Equal(t, req.GetSessionData(UserEmailSessionKey), nil, "User is logged in")
}

func TestLogoutViewsSuccess(t *testing.T) {
	helios.App.BeforeTest()

	requestData := LoginRequest{Email: "email", Password: "password"}
	sessionData := make(map[string]interface{})
	sessionData[UserEmailSessionKey] = "abc"

	req := helios.MockRequest{
		RequestData: requestData,
		SessionData: sessionData,
	}

	LogoutView(&req)

	assert.Equal(t, http.StatusOK, req.StatusCode, "Unexpected status code")
	assert.Empty(t, sessionData[UserEmailSessionKey], "User Email in session not changed")
}
