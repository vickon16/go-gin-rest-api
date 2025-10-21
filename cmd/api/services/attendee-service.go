package services

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vickon16/go-gin-rest-api/internal/app"
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

func CreateAttendee(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		var attendee models.Attendee

		if err := c.ShouldBindJSON(&attendee); err != nil {
			utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}

		if err := app.Models.Attendees.Insert(&attendee); err != nil {
			log.Printf("Error inserting attendee: %v", err)
			utils.ErrorResponse(c, "Failed to create attendee", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "Attendee Created successfully", models.CreateResponseAttendee(&attendee), http.StatusCreated)
	}
}

func GetAllAttendees(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {

		allAttendees, err := app.Models.Attendees.GetAll()
		if err != nil {
			log.Printf("Error getting attendees: %v", err)
			utils.ErrorResponse(c, "Failed to get attendees", http.StatusInternalServerError)
			return
		}
		if allAttendees == nil {
			utils.ErrorResponse(c, "Attendees not found", http.StatusNotFound)
			return
		}

		var serializedAttendees []models.AttendeeSerializer
		for _, attendee := range allAttendees {
			serializedAttendees = append(serializedAttendees, models.CreateResponseAttendee(attendee))
		}

		utils.SuccessResponse(c, "Successfully retrieved attendees", serializedAttendees)
	}
}

func GetEventsForAttendee(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		attendeeId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid attendee Id", http.StatusBadRequest)
			return
		}

		events, err := app.Models.Events.GetEventsByAttendeeId(attendeeId)
		if err != nil {
			log.Printf("Error getting events for attendee: %v", err)
			utils.ErrorResponse(c, "Failed to get events for attendee", http.StatusInternalServerError)
			return
		}
		if events == nil {
			utils.ErrorResponse(c, "No events found for attendee", http.StatusNotFound)
			return
		}

		var serializedEvents []models.EventSerializer
		for _, event := range events {
			serializedEvents = append(serializedEvents, models.CreateResponseEvent(event))
		}

		utils.SuccessResponse(c, "Successfully retrieved events for attendee", serializedEvents)
	}
}

func GetAttendee(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid attendee Id", http.StatusBadRequest)
			return
		}

		attendee, err := app.Models.Attendees.Get(id)
		if err != nil {
			utils.ErrorResponse(c, "Failed to get attendee", http.StatusInternalServerError)
			return
		}
		if attendee == nil {
			utils.ErrorResponse(c, "Attendee does not exist", http.StatusConflict)
			return
		}

		utils.SuccessResponse(c, "Successfully retrieved attendee", models.CreateResponseAttendee(attendee))
	}
}

func UpdateAttendee(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid attendee Id", http.StatusBadRequest)
			return
		}

		// Check for existing attendee
		existingAttendee, err := app.Models.Attendees.Get(id)
		if err != nil {
			utils.ErrorResponse(c, "Failed to get attendee", http.StatusInternalServerError)
			return
		}
		if existingAttendee == nil {
			utils.ErrorResponse(c, "No attendee found", http.StatusNotFound)
			return
		}

		var updatedAttendee models.UpdateAttendeeDto
		if err := c.ShouldBindJSON(&updatedAttendee); err != nil {
			utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}

		attendee, err := app.Models.Attendees.Update(id, &updatedAttendee)
		if err != nil {
			utils.ErrorResponse(c, "Failed to update attendee", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "Successfully updated attendee", models.CreateResponseAttendee(attendee))
	}
}

func DeleteAttendee(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid attendee Id", http.StatusBadRequest)
			return
		}

		if err := app.Models.Attendees.Delete(id); err != nil {
			utils.ErrorResponse(c, "Failed to delete attendee", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "Successfully deleted attendee", nil)
	}
}
