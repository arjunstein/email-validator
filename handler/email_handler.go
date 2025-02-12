package handler

import (
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
    Status string `json:"status"`
    Message string `json:"message"`
    Email string `json:"email"`
}

func CheckEmailHandler(c *gin.Context)  {
    godotenv.Load()

    var request EmailRequest

    // Bind JSON input
    if err := c.ShouldBindJSON(&request); err != nil {
        // when email is empty
        if request.Email == "" {
            c.JSON(http.StatusBadRequest, EmailResponse{
                Status: "error",
                Message: "Email is required",
                Email: "",
            })
            return
        }
    }

    // Validate email format
    if err := checkmail.ValidateFormat(request.Email); err != nil {
        c.JSON(http.StatusBadRequest, EmailResponse{
            Status: "error",
            Message: "Invalid email format",
            Email: request.Email,
        })
        return
    }

    // Validate MX record
    if err := checkmail.ValidateMX(request.Email); err != nil {
        c.JSON(http.StatusBadRequest, EmailResponse{
            Status: "error",
            Message: "Invalid domain email",
            Email: request.Email,
        })
        return
    }

    // Additional SMTP validation
    var (
        serverHostName    = os.Getenv("SMTP_HOST") 
        serverMailAddress = os.Getenv("SMTP_MAIL")
    )

    if err := checkmail.ValidateHostAndUser(serverHostName, serverMailAddress, request.Email); err != nil {
            c.JSON(http.StatusOK, EmailResponse{
                Status: "error",
                Message: "Email does not exist",
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