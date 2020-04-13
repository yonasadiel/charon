package auth

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeLoginRequest(t *testing.T) {
	var user User = UserFactory(User{Name: "User 1", Username: "user1", Password: "abcd"})
	var expectedJSON string = `{"name":"User 1","username":"user1","role":"participant"}`
	var serialized []byte
	var errMarshalling error

	serialized, errMarshalling = json.Marshal(SerializeUser(user))
	assert.Nil(t, errMarshalling, "Failed to marshalling user to json")
	assert.Equal(t, expectedJSON, string(serialized), "Different serialization")
}
