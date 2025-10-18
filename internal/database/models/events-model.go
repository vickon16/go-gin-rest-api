package models

import (
	"database/sql"
	"time"
)

type EventModel struct {
	DB *sql.DB
}

type Event struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId" binding:"required"`
	Name        string    `json:"name" binding:"required,min=3"`
	Description string    `json:"description" binding:"required,min=10,max=100"`
	Date        time.Time `json:"date" binding:"required"`
	Location    string    `json:"location" binding:"required,min=3"`
}

type UpdateEventDto struct {
	UserID      int       `json:"userId"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Location    string    `json:"location"`
}
