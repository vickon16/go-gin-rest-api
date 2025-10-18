package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/cmd/api/services"
	"github.com/vickon16/go-gin-rest-api/internal/app"
)

func setupAuthControllers(router *gin.RouterGroup, app *app.Application) {
	r := router.Group("/auth")

	r.POST("/", services.RegisterUser(app))
	r.GET("/", services.GetAllEvent(app))
	r.GET("/:id", services.GetEvent(app))
	r.PUT("/:id", services.UpdateEvent(app))
	r.DELETE("/:id", services.DeleteEvent(app))
}
