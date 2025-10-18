package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

func (app *application) createEvent(c *gin.Context) {
	var event models.Event
	log.Print("Here")

	if err := c.ShouldBindJSON(&event); err != nil {
		utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
	}

	log.Print(event)
}

func (app *application) getAllEvent(c *gin.Context) {
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {

	}
}

func (app *application) getEvent(c *gin.Context) {
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {

	}
}

func (app *application) updateEvent(c *gin.Context) {
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {

	}
}

func (app *application) deleteEvent(c *gin.Context) {
	var event models.Event

	if err := c.ShouldBindJSON(&event); err != nil {

	}
}
