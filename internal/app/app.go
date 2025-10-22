package app

import (
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
	"github.com/vickon16/go-gin-rest-api/internal/redisDb"
)

type Application struct {
	Port   int
	Models models.Models
	Redis  *redisDb.RedisClient
}
