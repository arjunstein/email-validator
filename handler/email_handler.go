package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
    SmtpMessage string `json:"smtp_message"`
}

func CheckEmailHandler(c *gin.Context)  {
    godotenv.Load()

    var request EmailRequest

    // Bind JSON input
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, EmailResponse{
            Status: "error",
            Message: "Invalid request",
            Email: "",
        })
        return
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
        if smtpErr, ok := err.(checkmail.SmtpError); ok {
            smtpMessage := fmt.Sprintf("SMTP Error - Code: %s, Msg: %s", smtpErr.Code(), smtpErr)

            // Log the SMTP message
            logDir := "logs"
            if _, err := os.Stat(logDir); os.IsNotExist(err) {
                os.Mkdir(logDir, 0755)
            }
            logFile, logErr := os.OpenFile(filepath.Join(logDir, "smtp_errors.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
            if logErr != nil {
                log.Println("Failed to open log file:", logErr)
            }
            logger := log.New(logFile, fmt.Sprintf("SMTP_ERROR [%s]: ", request.Email), log.LstdFlags)
            logger.Println(smtpMessage)

            c.JSON(http.StatusOK, EmailResponse{
                Status: "error",
                Message: "Email does not exist",
                SmtpMessage: smtpMessage,
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