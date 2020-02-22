package auth

import (
	"time"

	"github.com/yonasadiel/helios"
)

// User is that is registered in Charon app
type User struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"size:256"`
	Email    string `gorm:"size:256; unique"`
	Password string `gorm:"size:256"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func init() {
	helios.App.RegisterModel(User{})
}

// NewUser creates user with provided name, email, and password
// password will be hashed first
func NewUser(name, email, password string) User {
	return User{
		Name:     name,
		Email:    email,
		Password: hashPassword(password),
	}
}
