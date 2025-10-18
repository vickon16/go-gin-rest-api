package main

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
	"github.com/vickon16/go-gin-rest-api/internal/env"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	port      int
	jWTSecret string
	models    models.Models
}

func main() {
	_ = godotenv.Load()
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	models := models.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jWTSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}
}
