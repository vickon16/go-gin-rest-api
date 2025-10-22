package models

import (
	"database/sql"
	"time"
)

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	ID          int64     `db:"id" json:"id,omitempty"`
	UserID      int64     `db:"user_id" json:"userId,omitempty" binding:"required"`
	Name        string    `db:"name" json:"name,omitempty" binding:"required,min=3,max=255"`
	Description string    `db:"description" json:"description,omitempty" binding:"required,min=5"`
	Date        time.Time `db:"date" json:"date,omitempty" binding:"required"`
	Location    string    `db:"location" json:"location,omitempty" binding:"required"`
	BaseModel

	// Joins
	User *User `json:"user,omitempty"`
}

type CreateEventDto struct {
	UserID      int64     `json:"userId,omitempty"`
	Name        string    `json:"name,omitempty" binding:"required,min=3,max=255"`
	Description string    `json:"description,omitempty" binding:"required,min=5"`
	Date        time.Time `json:"date,omitempty" binding:"required"`
	Location    string    `json:"location,omitempty" binding:"required"`
}

type UpdateEventDto struct {
	UserID      int64     `json:"userId,omitempty"`
	Name        string    `json:"name,omitempty" binding:"omitempty,min=3,max=255"`
	Description string    `json:"description,omitempty" binding:"omitempty,min=5"`
	Date        time.Time `json:"date,omitempty" binding:"omitempty"`
	Location    string    `json:"location,omitempty" binding:"omitempty"`
}

type EventSerializer struct {
	ID          int64     `json:"id,omitempty"`
	UserID      int64     `json:"userId,omitempty"`
	Name        string    `json:"name,omitempty"`
	Description string    `json:"description,omitempty"`
	Date        time.Time `json:"date,omitempty"`
	Location    string    `json:"location,omitempty"`
	BaseModel

	// Joins
	User *UserSerializer `json:"user,omitempty"`
}

func CreateResponseEvent(event *Event) EventSerializer {
	response := EventSerializer{
		ID:          event.ID,
		UserID:      event.UserID,
		Name:        event.Name,
		Description: event.Description,
		Date:        event.Date,
		Location:    event.Location,
		BaseModel:   BaseModel{CreatedAt: event.CreatedAt, UpdatedAt: event.UpdatedAt},
	}

	if event.User != nil {
		userResponse := CreateResponseUser(event.User)
		response.User = &userResponse
	}

	return response
}
