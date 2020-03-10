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
	userLocal := User{Role: userRoleLocal}
	userParticipant := User{Role: userRoleParticipant}

	assert.True(t, userLocal.IsLocal(), "userLocal should be local")
	assert.False(t, userLocal.IsParticipant(), "userLocal should be local")
	assert.True(t, userParticipant.IsParticipant(), "userParticipant should be participant")
	assert.False(t, userLocal.IsParticipant(), "userLocal should be participant")

	userLocal.SetAsParticipant()
	assert.Equal(t, userRoleParticipant, userLocal.Role, "userLocal have been converted to participant")
	userLocal.SetAsLocal()
	assert.Equal(t, userRoleLocal, userLocal.Role, "userLocal have been converted to local")
}
