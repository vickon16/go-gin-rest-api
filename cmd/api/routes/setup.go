package routes

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vickon16/go-gin-rest-api/internal/app"
	"github.com/vickon16/go-gin-rest-api/internal/env"
	"github.com/vickon16/go-gin-rest-api/internal/utils"
)

func SetupRoutes(app *app.Application) http.Handler {
	g := gin.Default()

	// ✅ Custom Recovery Middleware
	g.Use(gin.Recovery())

	// ✅ CORS Middleware
	g.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000, https://yourfrontenddomain.com")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// Centralized error handling
	g.Use(func(c *gin.Context) {
		c.Next() // Process request

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			code := http.StatusInternalServerError
			message := err.Error()

			// if it's a Gin error with code, use it
			if c.Writer.Status() != http.StatusOK {
				code = c.Writer.Status()
			}

			// log internal errors
			if code == http.StatusInternalServerError {
				log.Printf("Internal Server Error: %v", err)
			}

			utils.ErrorResponse(c, message, code)
			c.Abort()
			return
		}
	})

	v1 := g.Group("/api/v1")

	// Auth
	setupAuthControllers(v1, app)

	// User
	setupUserControllers(v1, app)

	// Events
	setupEventsControllers(v1, app)

	// Attendees
	setupAttendeesControllers(v1, app)

	// for swagger docs
	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(302, "/swagger/index.html")
		}

		apiUrl := env.GetEnvString("API_URL", "http://localhost:8080")
		swaggerUrl := apiUrl + "/swagger/doc.json"

		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL(swaggerUrl))(c)
		// run 'swag init --dir cmd/api --parseDependency --parseInternal --parseDepth 1' to generate docs
	})

	return g
}
