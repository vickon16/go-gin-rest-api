package main

import (
	"database/sql"
	"log"

	"github.com/joho/godotenv"
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

	if err := serve(app); err != nil {
		log.Fatal(err)
	}
}
