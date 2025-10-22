package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/cmd/api/middlewares"
	"github.com/vickon16/go-gin-rest-api/cmd/api/services"
	"github.com/vickon16/go-gin-rest-api/internal/app"
)

func setupUserControllers(router *gin.RouterGroup, app *app.Application) {

	// Private routes
	priv := router.Group("/users", middlewares.AuthMiddleware(app))

	priv.GET("/", services.GetAllUsers(app))
	priv.GET("/:id", services.GetUser(app))
	priv.GET("/me", services.GetMe(app))
	priv.PUT("/:id", services.UpdateUser(app))
	priv.DELETE("/:id", services.DeleteUser(app))
}
