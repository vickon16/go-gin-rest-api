package models

import (
	"database/sql"
)

type AuthModel struct {
	DB *sql.DB
}

type RegisterUserDto struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

type LoginUserDto struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

type LoginSerializer struct {
	Token string `json:"token"`
}
