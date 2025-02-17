package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/arjunstein/email-validator/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main()  {
	// Load environment variable
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	r := gin.Default()

	apiPath := os.Getenv("PATH_URL")
	port := os.Getenv("APP_PORT")
	apiKey := os.Getenv("API_KEY")

	// Middleware to check API key
	r.Use(func(c *gin.Context) {
        clientAPIKey := c.GetHeader("x-api-key")
        if clientAPIKey != apiKey {
            c.JSON(http.StatusUnauthorized, gin.H{
                "status":  "error",
                "message": "Invalid API key",
            })
            c.Abort()
            return
        }
        c.Next()
    })

	// endpoint
	r.POST(apiPath, handler.CheckEmailHandler)

	r.Run(":" + port)
}