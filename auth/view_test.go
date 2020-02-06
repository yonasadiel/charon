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
}
