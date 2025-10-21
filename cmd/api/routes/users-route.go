package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/cmd/api/services"
	"github.com/vickon16/go-gin-rest-api/internal/app"
)

func setupAuthControllers(router *gin.RouterGroup, app *app.Application) {
	r := router.Group("/auth/users")

	r.POST("/", services.RegisterUser(app))
	r.GET("/", services.GetAllUsers(app))
	r.GET("/:id", services.GetUser(app))
	r.PUT("/:id", services.UpdateUser(app))
	r.DELETE("/:id", services.DeleteUser(app))
}
