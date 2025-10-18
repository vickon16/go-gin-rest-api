package models

import "database/sql"

type AttendeesModel struct {
	DB *sql.DB
}

type Attendees struct {
	ID      int `json:"id"`
	UserID  int `json:"userId" binding:"required"`
	EventID int `json:"eventId" binding:"required"`
}
