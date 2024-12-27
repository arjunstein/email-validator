package main

import (
	"fmt"
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

	// endpoint
	r.POST(apiPath, handler.CheckEmailHandler)

	r.Run(":" + port)
}