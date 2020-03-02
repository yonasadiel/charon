package auth

import (
	"time"

	"github.com/yonasadiel/helios"
)

// User is that is registered in Charon app. Types of user:
// - "local": user that organize of local exam.
// - "participant": user that taking the exam.
type User struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"size:256"`
	Email    string `gorm:"size:256; unique"`
	Password string `gorm:"size:256"`
	userType string `gorm:"size:10"` // enum("local", "participant"), default to "participant"

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Session of user logged in
type Session struct {
	ID        uint `gorm:"primary_key"`
	UserID    uint
	Token     string `gorm:"size:20;unique"`
	IPAddress string `gorm:"size:20"`

	User *User `gorm:"foreignkey:user_id"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func init() {
	helios.App.RegisterModel(User{})
	helios.App.RegisterModel(Session{})
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

// IsLocal returns true if the user is local
func (user User) IsLocal() bool {
	return user.userType == userTypeLocal
}

// IsParticipant returns true if the user is participant who taking the exam.
// This is the default value
func (user User) IsParticipant() bool {
	return user.userType != userTypeLocal
}
