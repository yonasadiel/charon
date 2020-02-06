package auth

import (
	"time"

	"github.com/yonasadiel/charon/app"
)

// User is that is registered in Charon app
type User struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name" validate:"required" gorm:"size:256"`
	Email    string `json:"email" validate:"required" gorm:"size:256; unique"`
	Password string `json:"-" gorm:"size:256"`

	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-"`
}

func init() {
	app.Charon.RegisterModel(User{})
}

// NewUser creates user with provided name, email, and password
// password will be hashed first
func NewUser(name, email, password string) (*User, error) {
	return &User{
		Name:     name,
		Email:    email,
		Password: hashPassword(password),
	}, nil
}
