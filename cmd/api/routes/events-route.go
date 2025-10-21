package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/cmd/api/services"
	"github.com/vickon16/go-gin-rest-api/internal/app"
)

func setupEventsControllers(router *gin.RouterGroup, app *app.Application) {
	r := router.Group("/events")

	r.POST("/", services.CreateEvent(app))
	r.GET("/", services.GetAllEvent(app))
	r.GET("/:id", services.GetEvent(app))
	r.POST("/:id/attendees/:userId", services.AddAttendeeToEvent(app))
	r.GET("/:id/attendees", services.GetAttendeesForEvent(app))
	r.PUT("/:id", services.UpdateEvent(app))
	r.DELETE("/:id", services.DeleteEvent(app))
}
