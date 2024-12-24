package handler

import (
	"net/http"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
)

type EmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type EmailResponse struct {
	Email string `json:"email"`
	Status string `json:"status"`
	Message string `json:"message"`
}

func CheckEmailHandler(c *gin.Context)  {
	var request EmailRequest

	// input validate
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"message": err.Error(),
		})
		return
	}

	// check email
	err := checkmail.ValidateFormat(request.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, EmailResponse{
			Status: "error",
			Message: "Invalid email format",
			Email: request.Email,
		})
		return
	}

	// check email connection
	err = checkmail.ValidateHost(request.Email)
	if err != nil {
		c.JSON(http.StatusOK, EmailResponse{
			Status: "error",
			Message: "Email not active or not exist",
			Email: request.Email,
		})
		return
	}

	// email valid
	c.JSON(http.StatusOK, EmailResponse{
		Status: "success",
		Message: "Email is valid",
		Email: request.Email,
	})
}