package main

import (
	"github.com/arjunstein/email-validator/handler"
	"github.com/gin-gonic/gin"
)

func main()  {
	r := gin.Default()

	// endpoint
	r.POST("/api/v1/check-email", handler.CheckEmailHandler)

	r.Run(":8080")
}