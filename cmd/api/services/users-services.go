package services

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/internal/app"
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

func RegisterUser(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto models.RegisterUserDto

		if err := c.ShouldBindJSON(&dto); err != nil {
			utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if user already exists
		if _, err := app.Models.Users.GetUserByEmail(dto.Email); err == nil {
			utils.ErrorResponse(c, "User already exists", http.StatusConflict)
			return
		}

		hashedPassword, err := utils.HashPassword(dto.Password)
		if err != nil {
			utils.ErrorResponse(c, "Something went wrong", http.StatusInternalServerError)
			return
		}

		user := models.User{
			Email:    dto.Email,
			Password: hashedPassword,
			Name:     dto.Name,
		}

		if err := app.Models.Users.Insert(&user); err != nil {
			utils.ErrorResponse(c, "Failed to create user", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "User Created successfully", user, http.StatusCreated)
	}
}
