package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/vickon16/go-gin-rest-api/cmd/api/routes"
	"github.com/vickon16/go-gin-rest-api/internal/app"
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
	"github.com/vickon16/go-gin-rest-api/internal/env"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	_ = godotenv.Load()
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	models := models.NewModels(db)
	app := &app.Application{
		Port:      env.GetEnvInt("PORT", 8080),
		JWTSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		Models:    models,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.Port),
		Handler:      routes.SetupRoutes(app),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Starting server on port %d", app.Port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}
