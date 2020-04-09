package exam

import (
	"time"

	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

// Event is the exam event
// It also stores the start and end time
type Event struct {
	ID          uint   `gorm:"primary_key"`
	Description string `gorm:"type:text"`
	Title       string `gorm:"size:100"`
	StartsAt    time.Time
	EndsAt      time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// UserEvent is many to many indicating an user is participating
// in an event.
type UserEvent struct {
	ID      uint `gorm:"primary_key"`
	EventID uint
	UserID  uint

	Event *Event     `gorm:"foreignkey:EventID"`
	User  *auth.User `gorm:"foreignkey:UserID"`

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
	ID         uint `gorm:"primary_key"`
	QuestionID uint
	UserID     uint
	Ordering   uint

	Question *Question  `gorm:"foreignkey:QuestionID"`
	User     *auth.User `gorm:"foreignkey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Submission of user to the question. This model should be insert only,
// never deleted, for log purpose.
type Submission struct {
	ID         uint   `gorm:"primary_key"`
	Answer     string `gorm:"type:text"`
	QuestionID uint
	UserID     uint

	Question *Question  `gorm:"foreignkey:QuestionID"`
	User     *auth.User `gorm:"foreignkey:UserID"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func init() {
	helios.App.RegisterModel(Event{})
	helios.App.RegisterModel(UserEvent{})
	helios.App.RegisterModel(QuestionChoice{})
	helios.App.RegisterModel(Question{})
	helios.App.RegisterModel(UserQuestion{})
	helios.App.RegisterModel(Submission{})
}
