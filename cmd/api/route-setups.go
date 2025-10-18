package main

import "github.com/gin-gonic/gin"

func (app *application) setupEventsRoutes(router *gin.RouterGroup) {
	r := router.Group("/events")

	r.POST("/", app.createEvent)
	r.GET("/", app.getAllEvent)
	r.GET("/:id", app.getEvent)
	r.PUT("/:id", app.updateEvent)
	r.DELETE("/:id", app.deleteEvent)
}
