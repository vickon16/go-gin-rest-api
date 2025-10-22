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

	User  *User  `json:"user,omitempty"`
	Event *Event `json:"event,omitempty"`
}

type CreateAttendeeDto struct {
	UserID  int64 `json:"userId" binding:"required"`
	EventID int64 `json:"eventId" binding:"required"`
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
	User  *UserSerializer  `json:"user,omitempty"`
	Event *EventSerializer `json:"event,omitempty"`
}

func CreateResponseAttendee(attendee *Attendee) AttendeeSerializer {

	response := AttendeeSerializer{
		ID:        attendee.ID,
		UserID:    attendee.UserID,
		EventID:   attendee.EventID,
		BaseModel: BaseModel{CreatedAt: attendee.CreatedAt, UpdatedAt: attendee.UpdatedAt},
	}

	if attendee.User != nil {
		userResponse := CreateResponseUser(attendee.User)
		response.User = &userResponse
	}

	if attendee.Event != nil {
		eventResponse := CreateResponseEvent(attendee.Event)
		response.Event = &eventResponse
	}

	return response
}
