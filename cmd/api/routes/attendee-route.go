package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/cmd/api/middlewares"
	"github.com/vickon16/go-gin-rest-api/cmd/api/services"
	"github.com/vickon16/go-gin-rest-api/internal/app"
)

func setupAttendeesControllers(router *gin.RouterGroup, app *app.Application) {
	priv := router.Group("/attendees", middlewares.AuthMiddleware(app))

	priv.POST("/", services.CreateAttendee(app))

	priv.GET("/", services.GetAllAttendees(app))
	priv.GET("/:id", services.GetAttendee(app))
	priv.GET("/:id/events", services.GetEventsByAttendee(app))

	priv.PUT("/:id", services.UpdateAttendee(app))

	priv.DELETE("/:id", services.DeleteAttendee(app))
}
