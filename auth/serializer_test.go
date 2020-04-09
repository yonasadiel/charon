package auth

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/yonasadiel/helios"
)

func TestSerializeLoginRequest(t *testing.T) {
	helios.App.BeforeTest()

	var user User = NewUser("User 1", "user1", "abcd")

	expected := `{"name":"User 1","username":"user1","role":"participant"}`
	actual, err := json.Marshal(SerializeUser(user))
	assert.Nil(t, err, "Failed to marshalling user to json")
	assert.Equal(t, expected, string(actual), "Different serialization")
}
