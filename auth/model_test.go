package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	user := NewUser("john", "johndoe@gmail.com", "password")

	assert.Equal(t, "john", user.Name, "Different user name")
	assert.Equal(t, "johndoe@gmail.com", user.Email, "Different user email")
	assert.True(t, checkPasswordHash("password", user.Password), "Different user password")
}

func TestNewUserErrorHash(t *testing.T) {
	user := NewUser("john", "johndoe@gmail.com", "password")

	assert.Equal(t, "john", user.Name, "Different user name")
	assert.Equal(t, "johndoe@gmail.com", user.Email, "Different user email")
	assert.True(t, checkPasswordHash("password", user.Password), "Different user password")
}

func TestUserType(t *testing.T) {
	user := NewUser("john", "johndoe@gmail.com", "password")

	assert.True(t, user.IsParticipant(), "user should be participant")
	assert.False(t, user.IsLocal(), "user should be participant")
	assert.False(t, user.IsOrganizer(), "user should be participant")
	assert.False(t, user.IsAdmin(), "user should be participant")

	user.SetAsLocal()
	assert.False(t, user.IsParticipant(), "user should be local")
	assert.True(t, user.IsLocal(), "user should be local")
	assert.False(t, user.IsOrganizer(), "user should be local")
	assert.False(t, user.IsAdmin(), "user should be local")

	user.SetAsOrganizer()
	assert.False(t, user.IsLocal(), "user should be organizer")
	assert.False(t, user.IsParticipant(), "user should be organizer")
	assert.True(t, user.IsOrganizer(), "user should be organizer")
	assert.False(t, user.IsAdmin(), "user should be organizer")

	user.SetAsAdmin()
	assert.False(t, user.IsLocal(), "user should be admin")
	assert.False(t, user.IsParticipant(), "user should be admin")
	assert.False(t, user.IsOrganizer(), "user should be admin")
	assert.True(t, user.IsAdmin(), "user should be admin")

	user.SetAsParticipant()
	assert.True(t, user.IsParticipant(), "user should be participant")
	assert.False(t, user.IsLocal(), "user should be participant")
	assert.False(t, user.IsOrganizer(), "user should be participant")
	assert.False(t, user.IsAdmin(), "user should be participant")
}
