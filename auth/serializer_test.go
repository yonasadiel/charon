package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yonasadiel/charon/app"
)

func TestDeserializeLoginRequest(t *testing.T) {
	app.Charon.BeforeTest()

	requestData := make(map[string]string)
	requestData["email"] = "abc"
	requestData["password"] = "def"
	requestData["other"] = "ghi"

	loginRequest := DeserializeLoginRequest(requestData)
	expectedLoginRequest := LoginRequest{Email: "abc", Password: "def"}
	assert.Equal(t, expectedLoginRequest, loginRequest, "Different deserialization")
}
