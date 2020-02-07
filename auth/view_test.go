package auth

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yonasadiel/charon/app"
)

func TestLoginViewsSuccess(t *testing.T) {
	app.Charon.BeforeTest()
	user, _ := NewUser("name", "email", "password")
	app.Charon.DB.Create(user)

	requestData := make(map[string]string)
	requestData["email"] = "email"
	requestData["password"] = "password"

	req := app.MockRequest{
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

	assert.Equal(t, user.ID, uint(returnedUser["id"].(float64)), "Wrong ID")
	assert.Equal(t, "name", returnedUser["name"], "Wrong Name")
	assert.Equal(t, "email", returnedUser["email"], "Wrong Email")
	assert.Equal(t, req.GetSessionData(UserIDSessionKey), user.ID, "Session is not changed")
}

func TestLoginViewWrongUsername(t *testing.T) {
	app.Charon.BeforeTest()
	user, _ := NewUser("name", "email", "password")
	app.Charon.DB.Create(user)

	requestData := make(map[string]string)
	requestData["email"] = "wrong_email"
	requestData["password"] = "password"

	req := app.MockRequest{
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
	assert.Equal(t, req.GetSessionData(UserIDSessionKey), nil, "User is logged in")
}

func TestLoginViewWrongPassword(t *testing.T) {
	app.Charon.BeforeTest()
	user, _ := NewUser("name", "email", "password")
	app.Charon.DB.Create(user)

	requestData := make(map[string]string)
	requestData["email"] = "email"
	requestData["password"] = "wrong_password"

	req := app.MockRequest{
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
	assert.Equal(t, req.GetSessionData(UserIDSessionKey), nil, "User is logged in")
}

func TestLogoutViewsSuccess(t *testing.T) {
	app.Charon.BeforeTest()

	requestData := make(map[string]string)
	requestData["email"] = "email"
	requestData["password"] = "password"

	sessionData := make(map[string]interface{})
	sessionData[UserIDSessionKey] = 2

	req := app.MockRequest{
		RequestData: requestData,
		SessionData: sessionData,
	}

	LogoutView(&req)

	assert.Equal(t, http.StatusOK, req.StatusCode, "Unexpected status code")
	assert.Equal(t, 0, sessionData[UserIDSessionKey], "User ID in session not changed")
}
