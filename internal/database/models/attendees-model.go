package models

import (
	"database/sql"
)

type AttendeesModel struct {
	DB *sql.DB
}

type Attendee struct {
	ID      int64 `db:"id" json:"id"`
	UserID  int64 `db:"user_id" json:"userId" binding:"required"`
	EventID int64 `db:"event_id" json:"eventId" binding:"required"`
	BaseModel

	User  *User
	Event *Event
}

type UpdateAttendeeDto struct {
	UserID  int64 `json:"userId"`
	EventID int64 `json:"eventId"`
}

type AttendeeSerializer struct {
	ID      int64 `json:"id,omitempty"`
	UserID  int64 `json:"userId,omitempty"`
	EventID int64 `json:"eventId,omitempty"`
	BaseModel

	// Joins
	User  UserSerializer  `json:"user"`
	Event EventSerializer `json:"event"`
}

func CreateResponseAttendee(attendee *Attendee) AttendeeSerializer {
	return AttendeeSerializer{
		ID:        attendee.ID,
		UserID:    attendee.UserID,
		EventID:   attendee.EventID,
		BaseModel: BaseModel{CreatedAt: attendee.CreatedAt, UpdatedAt: attendee.UpdatedAt},

		// Joins
		User:  CreateResponseUser(attendee.User),
		Event: CreateResponseEvent(attendee.Event),
	}
}
