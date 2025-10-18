package models

import "database/sql"

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=3"`
	Password string `json:"-" binding:"required,min=6,max=64"` // hidden in JSON responses
}
