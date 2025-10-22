package services

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/cmd/api/middlewares"
	"github.com/vickon16/go-gin-rest-api/internal/app"
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

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
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

func GetMe(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		contextUser := middlewares.GetUserFromContext(c)

		user, err := app.Models.Users.Get(contextUser.ID)
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
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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

		cacheKey := utils.ConstructRedisUserKey(id)
		app.Redis.Delete(cacheKey)

		utils.SuccessResponse(c, "Successfully updated user", models.CreateResponseUser(user))
	}
}

func DeleteUser(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
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
			utils.ErrorResponse(c, "User does not exist", http.StatusBadRequest)
			return
		}

		if err := app.Models.Users.Delete(id); err != nil {
			utils.ErrorResponse(c, "Failed to delete user", http.StatusInternalServerError)
			return
		}

		cacheKey := utils.ConstructRedisUserKey(id)
		app.Redis.Delete(cacheKey)

		utils.SuccessResponse(c, "Successfully deleted user", nil)
	}
}
