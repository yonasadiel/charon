package exam

import (
	"time"

	"github.com/yonasadiel/charon/backend/auth"
	"github.com/yonasadiel/helios"
)

// Event is the exam event
// It also stores the start and end time
type Event struct {
	ID          uint   `gorm:"primary_key"`
	Slug        string `gorm:"size:100;unique"`
	Title       string `gorm:"size:256"`
	Description string `gorm:"type:text"`
	Key         string `gorm:"size:48"`
	StartsAt    time.Time
	EndsAt      time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Venue is the event venue
type Venue struct {
	ID   uint `gorm:"primary_key"`
	Name string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Participation is many to many indicating an user is participating
// in a local event.
type Participation struct {
	ID      uint `gorm:"primary_key"`
	EventID uint
	UserID  uint
	VenueID uint

	Event *Event     `gorm:"foreignkey:EventID"`
	User  *auth.User `gorm:"foreignkey:UserID"`
	Venue *Venue     `gorm:"foreignkey:VenueID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Question that will be served to participant
// UserAnswer is the user current answer to related question
type Question struct {
	ID         uint   `gorm:"primary_key"`
	Content    string `gorm:"type:text"`
	EventID    uint
	UserAnswer string `gorm:"-"`

	Event   *Event           `gorm:"foreignkey:EventID"`
	Choices []QuestionChoice `gorm:"foreignkey:QuestionID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// QuestionChoice is choice of question with multiple choices type
type QuestionChoice struct {
	ID         uint `gorm:"primray_key"`
	Text       string
	QuestionID uint

	Question *Question `gorm:"foreignkey:QuestionID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// UserQuestion is the many-to-many relation between user and question
// because user can't see al questions, but only theirs.
type UserQuestion struct {
	ID              uint `gorm:"primary_key"`
	ParticipationID uint
	QuestionID      uint
	Ordering        uint
	Answer          string `gorm:"type:text"`

	Participation *Participation `gorm:"foreignkey:ParticipationID"`
	Question      *Question      `gorm:"foreignkey:QuestionID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func init() {
	helios.App.RegisterModel(Event{})
	helios.App.RegisterModel(Venue{})
	helios.App.RegisterModel(Participation{})
	helios.App.RegisterModel(QuestionChoice{})
	helios.App.RegisterModel(Question{})
	helios.App.RegisterModel(UserQuestion{})
}
