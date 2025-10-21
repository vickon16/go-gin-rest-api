package database

import (
	"database/sql"
	"log"

	"github.com/vickon16/go-gin-rest-api/internal/env"
)

func SetupDatabase() *sql.DB {
	dsn := env.GetEnvString("DATABASE_URL", "postgres://postgres:password@localhost:5432/mydb?sslmode=disable")
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatalf("Could not connect to PostgreSQL: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL database")

	return db
}

// Run this
// migrate -path ./cmd/migrate/migrations -database "postgres://postgres:Password123!@localhost:5432/go-gin-tutorial?sslmode=disable" up

// migrate -path ./cmd/migrate/migrations -database "postgres://postgres:Password123!@localhost:5432/go-gin-tutorial?sslmode=disable" force 1
