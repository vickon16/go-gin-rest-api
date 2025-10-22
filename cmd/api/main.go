package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/vickon16/go-gin-rest-api/cmd/api/routes"
	"github.com/vickon16/go-gin-rest-api/internal/app"
	"github.com/vickon16/go-gin-rest-api/internal/database"
	"github.com/vickon16/go-gin-rest-api/internal/database/models"
	"github.com/vickon16/go-gin-rest-api/internal/env"
	"github.com/vickon16/go-gin-rest-api/internal/redisDb"

	_ "github.com/lib/pq"
	_ "github.com/vickon16/go-gin-rest-api/docs"
)

// @title Go Gin Rest API
// @version 1.0
// @description A rest API in Go using Gin Framework
// @termsOfService  http://example.com/terms/

// @contact.name   Victor Cyril
// @contact.url    https://victorcyril.com
// @contact.email  vicy@victorcyril.com

// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter you bearer token in the format **Bearer &lt;token&gt;**

func main() {
	_ = godotenv.Load()

	db := database.SetupDatabase()
	defer db.Close()

	models := models.NewModels(db)
	redisClient := redisDb.NewRedisClient()

	app := &app.Application{
		Port:   env.GetEnvInt("PORT", 8080),
		Models: models,
		Redis:  redisClient,
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
