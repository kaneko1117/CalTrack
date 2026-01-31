package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func RespondError(c *gin.Context, status int, code, message string) {
	c.JSON(status, ErrorResponse{Code: code, Message: message})
}

func RespondValidationError(c *gin.Context, details []string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    CodeValidationError,
		Message: "Validation failed",
		Details: details,
	})
}
