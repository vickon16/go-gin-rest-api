package services

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/internal/app"
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

func CreateEvent(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		var event models.Event

		if err := c.ShouldBindJSON(&event); err != nil {
			utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}

		if err := app.Models.Events.Insert(&event); err != nil {
			utils.ErrorResponse(c, "Failed to create event", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "Event Created successfully", event, http.StatusCreated)
	}
}

func GetAllEvent(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {

		allEvents, err := app.Models.Events.GetAll()
		if err != nil {
			utils.ErrorResponse(c, "Failed to get events", http.StatusInternalServerError)
			return
		}
		if allEvents == nil {
			utils.ErrorResponse(c, "Events not found", http.StatusNotFound)
			return
		}

		utils.SuccessResponse(c, "Successfully retrieved events", allEvents)
	}
}

func GetEvent(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid event Id", http.StatusBadRequest)
			return
		}

		event, err := app.Models.Events.Get(id)
		if err != nil {
			utils.ErrorResponse(c, "Failed to get event", http.StatusInternalServerError)
			return
		}
		if event == nil {
			utils.ErrorResponse(c, "Event not found", http.StatusNotFound)
			return
		}

		utils.SuccessResponse(c, "Successfully retrieved event", event)
	}
}

func UpdateEvent(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid event Id", http.StatusBadRequest)
			return
		}

		// Check for existing event
		existingEvent, err := app.Models.Events.Get(id)
		if err != nil {
			utils.ErrorResponse(c, "Failed to get event", http.StatusInternalServerError)
			return
		}
		if existingEvent == nil {
			utils.ErrorResponse(c, "Event not found", http.StatusNotFound)
			return
		}

		var updatedEvent models.UpdateEventDto
		if err := c.ShouldBindJSON(&updatedEvent); err != nil {
			utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}

		if err := app.Models.Events.Update(id, &updatedEvent); err != nil {
			utils.ErrorResponse(c, "Failed to update event", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "Successfully updated event", updatedEvent)
	}
}

func DeleteEvent(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid event Id", http.StatusBadRequest)
			return
		}

		if err := app.Models.Events.Delete(id); err != nil {
			utils.ErrorResponse(c, "Failed to delete event", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "Successfully deleted event", nil)
	}
}
