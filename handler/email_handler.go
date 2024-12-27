package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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
	godotenv.Load()

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

	// Additional SMTP validation
	var (
		serverHostName    = os.Getenv("SMTP_HOST") 
		serverMailAddress = os.Getenv("SMTP_MAIL")
	)

	smtpErr := checkmail.ValidateHostAndUser(serverHostName, serverMailAddress, request.Email)
	if smtpErr != nil {
		if smtpErr, ok := smtpErr.(checkmail.SmtpError); ok {
			c.JSON(http.StatusOK, EmailResponse{
				Status: "error",
				Message: fmt.Sprintf("SMTP Error - Code: %s, Msg: %s", smtpErr.Code(), smtpErr),
				Email: request.Email,
			})
			return
		}
	}

	// email valid
	c.JSON(http.StatusOK, EmailResponse{
		Status: "success",
		Message: "Email is valid",
		Email: request.Email,
	})
}