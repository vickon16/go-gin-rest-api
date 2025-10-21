package services

import (
	"log"
	"net/http"
	"strconv"
	"sync"

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
			log.Printf("Error inserting event: %v", err)
			utils.ErrorResponse(c, "Failed to create event", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "Event Created successfully", models.CreateResponseEvent(&event), http.StatusCreated)
	}
}

func GetAllEvent(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {

		allEvents, err := app.Models.Events.GetAll()
		if err != nil {
			log.Printf("Error getting events: %v", err)
			utils.ErrorResponse(c, "Failed to get events", http.StatusInternalServerError)
			return
		}
		if allEvents == nil {
			utils.ErrorResponse(c, "Events not found", http.StatusNotFound)
			return
		}

		var serializedEvents []models.EventSerializer
		for _, event := range allEvents {
			serializedEvents = append(serializedEvents, models.CreateResponseEvent(event))
		}

		utils.SuccessResponse(c, "Successfully retrieved events", serializedEvents)
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
			utils.ErrorResponse(c, "Event does not exist", http.StatusConflict)
			return
		}

		utils.SuccessResponse(c, "Successfully retrieved event", models.CreateResponseEvent(event))
	}
}

func AddAttendeeToEvent(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventId, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid event Id", http.StatusBadRequest)
			return
		}
		userId, err := strconv.Atoi(c.Param("userId"))
		if err != nil {
			utils.ErrorResponse(c, "Invalid event Id", http.StatusBadRequest)
			return
		}

		var (
			event             *models.Event
			userToAdd         *models.User
			eventErr, userErr error
		)

		var wg sync.WaitGroup
		wg.Add(2)

		// Run event fetch in a goroutine
		go func() {
			defer wg.Done()
			event, eventErr = app.Models.Events.Get(eventId)
		}()

		// Run user fetch in a goroutine
		go func() {
			defer wg.Done()
			userToAdd, userErr = app.Models.Users.Get(userId)
		}()

		// Wait for both to finish
		wg.Wait()

		// Handle errors after both requests complete
		if eventErr != nil {
			log.Printf("Error getting event: %v", eventErr)
			utils.ErrorResponse(c, "Failed to get event", http.StatusNotFound)
			return
		}
		if userErr != nil {
			log.Printf("Error getting user: %v", userErr)
			utils.ErrorResponse(c, "Failed to get user", http.StatusNotFound)
			return
		}

		// Check if the user is not already an attendee
		existingAttendee, err := app.Models.Attendees.GetByEventAndAttendee(event.ID, userToAdd.ID)
		if err != nil {
			utils.ErrorResponse(c, "Failed to get events by attendee", http.StatusConflict)
			return
		}
		if existingAttendee != nil {
			utils.ErrorResponse(c, "Attendee already exist for this event", http.StatusConflict)
			return
		}

		attendee := models.Attendee{
			EventID: event.ID,
			UserID:  event.UserID,
		}

		err = app.Models.Attendees.Insert(&attendee)
		if err != nil {
			log.Printf("Error adding attendee: %v", userErr)
			utils.ErrorResponse(c, "Failed to add attendee to event", http.StatusConflict)
			return
		}

		utils.SuccessResponse(c, "Successfully added attendee to event", models.CreateResponseAttendee(&attendee), http.StatusCreated)
	}
}

func GetAttendeesForEvent(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		eventId, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			utils.ErrorResponse(c, "Invalid event Id", http.StatusBadRequest)
			return
		}

		attendees, err := app.Models.Attendees.GetAttendeesByEventId(eventId)
		if err != nil {
			log.Printf("Error getting attendees for event: %v", err)
			utils.ErrorResponse(c, "Failed to get attendees for event", http.StatusInternalServerError)
			return
		}

		var serializedAttendees []models.AttendeeSerializer
		for _, attendee := range attendees {
			serializedAttendees = append(serializedAttendees, models.CreateResponseAttendee(attendee))
		}

		utils.SuccessResponse(c, "Successfully retrieved attendees for event", serializedAttendees)
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
			utils.ErrorResponse(c, "No event found", http.StatusNotFound)
			return
		}

		var updatedEvent models.UpdateEventDto
		if err := c.ShouldBindJSON(&updatedEvent); err != nil {
			utils.ErrorResponse(c, err.Error(), http.StatusBadRequest)
			return
		}

		event, err := app.Models.Events.Update(id, &updatedEvent)
		if err != nil {
			utils.ErrorResponse(c, "Failed to update event", http.StatusInternalServerError)
			return
		}

		utils.SuccessResponse(c, "Successfully updated event", models.CreateResponseEvent(event))
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
