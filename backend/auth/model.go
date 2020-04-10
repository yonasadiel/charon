package auth

import (
	"time"

	"github.com/yonasadiel/helios"
)

// User is that is registered in Charon app. Types of user:
// - "local": user that organize of local exam.
// - "participant": user that taking the exam.
// - "admin": administrator of applicaton.
// - "organizer": writer of problems, etc.
type User struct {
	ID       uint   `gorm:"primary_key"`
	Name     string `gorm:"size:256"`
	Username string `gorm:"size:256; unique"`
	Password string `gorm:"size:256"`
	Role     string `gorm:"size:10"` // enum("local", "participant", "admin", "organizer"), default to "participant"

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

// NewUser creates user with provided name, username, and password
// password will be hashed first
func NewUser(name, username, password string) User {
	return User{
		Name:     name,
		Username: username,
		Password: hashPassword(password),
		Role:     userRoleParticipant,
	}
}

// IsAdmin returns true if the user is local
func (user *User) IsAdmin() bool {
	return user.Role == userRoleAdmin
}

// IsOrganizer returns true if the user is local
func (user *User) IsOrganizer() bool {
	return user.Role == userRoleOrganizer
}

// IsLocal returns true if the user is local
func (user *User) IsLocal() bool {
	return user.Role == userRoleLocal
}

// IsParticipant returns true if the user is participant who taking the exam.
// This is the default value
func (user *User) IsParticipant() bool {
	return user.Role == userRoleParticipant
}

// SetAsAdmin set the user as local administrator of exam
func (user *User) SetAsAdmin() {
	user.Role = userRoleAdmin
}

// SetAsOrganizer set the user as participant of exam
func (user *User) SetAsOrganizer() {
	user.Role = userRoleOrganizer
}

// SetAsLocal set the user as local administrator of exam
func (user *User) SetAsLocal() {
	user.Role = userRoleLocal
}

// SetAsParticipant set the user as participant of exam
func (user *User) SetAsParticipant() {
	user.Role = userRoleParticipant
}
