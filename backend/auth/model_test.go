package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	user1 := UserFactory(User{Password: "password"})
	user2 := User{Password: "password"}
	assert.True(t, checkPasswordHash("password", user1.Password))
	assert.False(t, checkPasswordHash("password", user2.Password))
}

func TestUserType(t *testing.T) {
	user := UserFactory(User{})

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
