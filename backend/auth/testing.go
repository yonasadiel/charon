package auth

import (
	"fmt"

	"github.com/yonasadiel/helios"
)

var userSeq uint = 0

// UserFactory creates an user for testing. The given argument will be
// completed if the attribute is empty.
func UserFactory(user User) User {
	userSeq = userSeq + 1
	if user.Name == "" {
		user.Name = fmt.Sprintf("Name User #%d", userSeq)
	}
	if user.Username == "" {
		user.Username = fmt.Sprintf("username%d", userSeq)
	}
	if user.Password == "" {
		user.Password = hashPassword(fmt.Sprintf("password-%d", userSeq))
	} else {
		user.Password = hashPassword(user.Password)
	}
	if user.Role == 0 {
		user.Role = UserRoleParticipant
	}
	return user
}

// UserFactorySaved do exactly like UserFactory but the result
// will be saved to database
func UserFactorySaved(user User) User {
	if user.ID == 0 {
		user = UserFactory(user)
		helios.DB.Create(&user)
	}
	return user
}
