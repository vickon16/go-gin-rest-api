package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/cmd/api/middlewares"
	"github.com/vickon16/go-gin-rest-api/cmd/api/services"
	"github.com/vickon16/go-gin-rest-api/internal/app"
)

func setupEventsControllers(router *gin.RouterGroup, app *app.Application) {
	priv := router.Group("/events", middlewares.AuthMiddleware(app))

	priv.POST("/", services.CreateEvent(app))
	priv.POST("/:id/attendees/:userId", services.AddAttendeeToEvent(app))

	priv.GET("/", services.GetAllEvent(app))
	priv.GET("/:id", services.GetEvent(app))
	priv.GET("/:id/attendees", services.GetAttendeesForEvent(app))

	priv.PUT("/:id", services.UpdateEvent(app))

	priv.DELETE("/:id/attendees/:userId", services.DeleteAttendeeFromEvent(app))
	priv.DELETE("/:id", services.DeleteEvent(app))
}
