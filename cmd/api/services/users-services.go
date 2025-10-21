package services

import (
	"log"
	"net/http"
	"strconv"

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

func GetAllUsers(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {

		allUsers, err := app.Models.Users.GetAll()
		if err != nil {
			utils.ErrorResponse(c, "Failed to get users", http.StatusInternalServerError)
			return
		}
		if allUsers == nil {
			utils.ErrorResponse(c, "Users not found", http.StatusNotFound)
			return
		}

		var serializedUsers []models.UserSerializer
		for _, user := range allUsers {
			serializedUsers = append(serializedUsers, models.CreateResponseUser(user))
		}

		utils.SuccessResponse(c, "Successfully retrieved users", serializedUsers)
	}
}

func GetUser(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid user Id", http.StatusBadRequest)
			return
		}

		user, err := app.Models.Users.Get(id)
		if err != nil {
			utils.ErrorResponse(c, "Failed to get user", http.StatusInternalServerError)
			return
		}
		if user == nil {
			utils.ErrorResponse(c, "User not found", http.StatusNotFound)
			return
		}

		utils.SuccessResponse(c, "Successfully retrieved user", models.CreateResponseUser(user))
	}
}

func UpdateUser(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid user Id", http.StatusBadRequest)
			return
		}

		// Check for existing user
		existingUser, err := app.Models.Users.Get(id)
		if err != nil {
			utils.ErrorResponse(c, "Failed to get user", http.StatusInternalServerError)
			return
		}
		if existingUser == nil {
			utils.ErrorResponse(c, "User not found", http.StatusNotFound)
			return
		}

		var updatedUser models.UpdateUserDto
		if err := c.ShouldBindJSON(&updatedUser); err != nil {
			utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := app.Models.Users.Update(id, &updatedUser)
		if err != nil {
			utils.ErrorResponse(c, "Failed to update user", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "Successfully updated user", models.CreateResponseUser(user))
	}
}

func DeleteUser(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid user Id", http.StatusBadRequest)
			return
		}

		if err := app.Models.Users.Delete(id); err != nil {
			utils.ErrorResponse(c, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "Successfully deleted user", nil)
	}
}
