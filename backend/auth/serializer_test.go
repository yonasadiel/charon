package auth

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerializeLoginRequest(t *testing.T) {
	var user User = UserFactory(User{Name: "User 1", Username: "user1", Password: "abcd"})

	expected := `{"name":"User 1","username":"user1","role":"participant"}`
	actual, err := json.Marshal(SerializeUser(user))
	assert.Nil(t, err, "Failed to marshalling user to json")
	assert.Equal(t, expected, string(actual), "Different serialization")
}
