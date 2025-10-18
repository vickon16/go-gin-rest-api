package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/vickon16/go-gin-rest-api/internal/app"
)

func serve(app *app.Application) error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Port),
		Handler:      routes(app),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting server on port %d", app.Port)
	return server.ListenAndServe()
}
