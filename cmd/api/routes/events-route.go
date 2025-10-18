package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/cmd/api/services"
	"github.com/vickon16/go-gin-rest-api/internal/app"
)

func SetupEventsRoutes(router *gin.RouterGroup, app *app.Application) {
	r := router.Group("/events")

	r.POST("/", services.CreateEvent())
	r.GET("/", services.GetAllEvents())
	r.GET("/:id", services.GetEvent())
	r.PUT("/:id", services.UpdateEvent())
	r.DELETE("/:id", services.DeleteEvent())
}
