package models

import (
	"database/sql"
)

type UserModel struct {
	DB *sql.DB
}

type User struct {
	ID       int64  `db:"id" json:"id"`
	Email    string `db:"email" json:"email" binding:"required,email"`
	Name     string `db:"name" json:"name" binding:"required,min=2,max=100"`
	Password string `db:"password" json:"-" binding:"required,min=6"`
	BaseModel
}

type RegisterUserDto struct {
	Email    string `json:"email" binding:"required,email"`
	Name     string `json:"name" binding:"required,min=3"`
	Password string `json:"password" binding:"required,min=6,max=64"`
}

type UpdateUserDto struct {
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Name     string `json:"name,omitempty" binding:"omitempty,min=3"`
	Password string `json:"password,omitempty" binding:"omitempty,min=6,max=64"`
}

type UserSerializer struct {
	ID    int64  `json:"id,omitempty"`
	Email string `json:"email,omitempty"`
	Name  string `json:"name,omitempty"`
	BaseModel
}

func CreateResponseUser(user *User) UserSerializer {
	return UserSerializer{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
		BaseModel: BaseModel{
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}
}
