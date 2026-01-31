package common

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func RespondError(c echo.Context, status int, code, message string) error {
	return c.JSON(status, ErrorResponse{Code: code, Message: message})
}

func RespondValidationError(c echo.Context, details []string) error {
	return c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    CodeValidationError,
		Message: "Validation failed",
		Details: details,
	})
}
