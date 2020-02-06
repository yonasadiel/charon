package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user, err := NewUser("john", "johndoe@gmail.com", "password")
	if err != nil {
		t.Fatalf("Fail on creating new user: %s", err)
	}

	assert.Equal(t, "john", user.Name, "Different user name")
	assert.Equal(t, "johndoe@gmail.com", user.Email, "Different user email")
	assert.True(t, checkPasswordHash("password", user.Password), "Different user password")
}

func TestNewUserErrorHash(t *testing.T) {
	user, err := NewUser("john", "johndoe@gmail.com", "password")
	if err != nil {
		t.Fatalf("Fail on creating new user: %s", err)
	}

	assert.Equal(t, "john", user.Name, "Different user name")
	assert.Equal(t, "johndoe@gmail.com", user.Email, "Different user email")
	assert.True(t, checkPasswordHash("password", user.Password), "Different user password")
}
