package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
)

func GetUserFromContext(ctx *gin.Context) *models.UserSerializer {
	contextUser, exists := ctx.Get("user")
	if !exists {
		return &models.UserSerializer{}
	}

	user, ok := contextUser.(*models.UserSerializer)
	if !ok {
		return &models.UserSerializer{}
	}

	return user
}
