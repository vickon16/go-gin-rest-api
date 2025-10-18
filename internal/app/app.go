package app

import "github.com/vickon16/go-gin-rest-api/internal/database/models"

type Application struct {
	Port      int
	JWTSecret string
	Models    models.Models
}
