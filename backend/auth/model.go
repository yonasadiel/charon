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
	ID            uint   `gorm:"primary_key"`
	Name          string `gorm:"size:256"`
	Username      string `gorm:"size:256; unique"`
	Password      string `gorm:"size:256"`
	Role          uint   // default to participant
	SessionLocked bool   `gorm:"default:false"`

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

// IsAdmin returns true if the user is local
func (user *User) IsAdmin() bool {
	return user.Role == UserRoleAdmin
}

// IsOrganizer returns true if the user is local
func (user *User) IsOrganizer() bool {
	return user.Role == UserRoleOrganizer
}

// IsLocal returns true if the user is local
func (user *User) IsLocal() bool {
	return user.Role == UserRoleLocal
}

// IsParticipant returns true if the user is participant who taking the exam.
// This is the default value
func (user *User) IsParticipant() bool {
	return user.Role == UserRoleParticipant
}

// SetAsAdmin set the user as local administrator of exam
func (user *User) SetAsAdmin() {
	user.Role = UserRoleAdmin
}

// SetAsOrganizer set the user as participant of exam
func (user *User) SetAsOrganizer() {
	user.Role = UserRoleOrganizer
}

// SetAsLocal set the user as local administrator of exam
func (user *User) SetAsLocal() {
	user.Role = UserRoleLocal
}

// SetAsParticipant set the user as participant of exam
func (user *User) SetAsParticipant() {
	user.Role = UserRoleParticipant
}
