package services

import (
	"log"
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
		existingUser, err := app.Models.Users.GetUserByEmail(dto.Email)
		if err != nil {
			utils.ErrorResponse(c, "Failed to get user by email", http.StatusInternalServerError)
			return
		}
		if existingUser != nil {
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
			log.Printf("Failed to create user %v", err)
			utils.ErrorResponse(c, "Failed to create user", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "User Created successfully", models.CreateResponseUser(&user), http.StatusCreated)
	}
}

func LoginUser(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dto models.LoginUserDto

		if err := c.ShouldBindJSON(&dto); err != nil {
			utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}

		// Check if user already exists
		existingUser, err := app.Models.Users.GetUserByEmail(dto.Email, true)
		if err != nil {
			utils.ErrorResponse(c, "Failed to get user by email", http.StatusInternalServerError)
			return
		}
		if existingUser == nil {
			utils.ErrorResponse(c, "User does not exists", http.StatusBadRequest)
			return
		}

		isValid := utils.CheckPasswordHash(existingUser.Password, dto.Password)
		if !isValid {
			utils.ErrorResponse(c, "Invalid email or password", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateLoginToken(existingUser.ID, existingUser.Email)
		if err != nil {
			utils.ErrorResponse(c, "Something went wrong", http.StatusInternalServerError)
			return
		}

		loginResponse := models.LoginSerializer{
			Token: token,
		}

		utils.SuccessResponse(c, "Login successful", loginResponse)
	}
}
