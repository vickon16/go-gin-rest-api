package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/internal/app"
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

func AuthMiddleware(app *app.Application) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(ctx, "Missing authorized header", http.StatusUnauthorized)
			// Cancel any other middleware in the chain
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			utils.ErrorResponse(ctx, "Bearer token is required", http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		// Parse and validate token
		claims, err := utils.VerifyToken(tokenString)
		if err != nil {
			utils.ErrorResponse(ctx, "Invalid or expired token", http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		userId := claims.UserID
		user, err := storeOrRetrieveFromRedis(app, userId)
		if err != nil {
			utils.ErrorResponse(ctx, err.Error(), http.StatusUnauthorized)
			ctx.Abort()
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

func storeOrRetrieveFromRedis(app *app.Application, userId int64) (*models.UserSerializer, error) {
	// find user
	cacheKey := utils.ConstructRedisUserKey(userId)
	cachedUser, err := app.Redis.Get(cacheKey)

	var user models.UserSerializer

	if err == nil && cachedUser != "" {
		// found in redis
		// Redis returns a string, convert back to json
		err := json.Unmarshal([]byte(cachedUser), &user)
		if err != nil {
			return nil, errors.New("something went wrong: 1")
		}
	} else {
		// Not in cache. Fetch from Db
		newUser, err := app.Models.Users.Get(userId)
		if err != nil || newUser == nil {
			return nil, errors.New("unauthorized user")
		}

		user = models.CreateResponseUser(newUser)

		// Store in redis
		userJson, err := json.Marshal(user)
		if err != nil {
			return nil, errors.New("something went wrong: 2")
		}

		// Set the json in redis
		err = app.Redis.Set(cacheKey, userJson, 30*time.Minute)
		if err != nil {
			return nil, errors.New("something went wrong: 3")
		}
	}

	return &user, nil
}
