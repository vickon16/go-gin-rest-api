package services

import (
	"errors"
	"log"
	"sync"

	"github.com/vickon16/go-gin-rest-api/internal/app"
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
)

func FindEventAndUser(app *app.Application, eventId, userId int64) (*models.Event, *models.User, error) {
	// Find existing attendee and event
	var (
		existingEvent     *models.Event
		existingUser      *models.User
		eventErr, userErr error
	)

	var wg sync.WaitGroup
	wg.Add(2)

	// Run event fetch in a goroutine
	go func() {
		defer wg.Done()
		existingEvent, eventErr = app.Models.Events.Get(eventId)
	}()

	// Run user fetch in a goroutine
	go func() {
		defer wg.Done()
		existingUser, userErr = app.Models.Users.Get(userId)
	}()

	// Wait for both to finish
	wg.Wait()

	// Handle errors after both requests complete
	if eventErr != nil {
		log.Printf("Error getting event: %v", eventErr)
		return nil, nil, errors.New("failed to get event")
	}
	if userErr != nil {
		log.Printf("Error getting user: %v", userErr)
		return nil, nil, errors.New("failed to get user")
	}

	if existingEvent == nil {
		return nil, nil, errors.New("event not found")
	}

	if existingUser == nil {
		return nil, nil, errors.New("user not found")
	}

	return existingEvent, existingUser, nil
}
