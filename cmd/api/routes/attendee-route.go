package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/cmd/api/services"
	"github.com/vickon16/go-gin-rest-api/internal/app"
)

func setupAttendeesControllers(router *gin.RouterGroup, app *app.Application) {
	r := router.Group("/attendees")

	r.POST("/", services.CreateAttendee(app))
	r.GET("/", services.GetAllAttendees(app))
	r.GET("/:id", services.GetAttendee(app))
	r.GET("/:id/events", services.GetEventsForAttendee(app))
	r.PUT("/:id", services.UpdateAttendee(app))
	r.DELETE("/:id", services.DeleteAttendee(app))
}
