package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/internal/app"
)

func routes(app *app.Application) http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")

	v1.POST("/events")

	// Events
	routes.SetupEventsRoutes(v1, app)

	return g
}
