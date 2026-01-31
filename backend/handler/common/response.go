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

// RespondError はエラーレスポンスを返す
func RespondError(c *gin.Context, status int, code, message string, originalErr error) {
	// 500系エラーの場合は詳細をログ出力
	if status >= 500 && originalErr != nil {
		LogError("RespondError", originalErr,
			"code", code,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
		)
	}
	c.JSON(status, ErrorResponse{Code: code, Message: message})
}

// RespondValidationError はバリデーションエラーレスポンスを返す
func RespondValidationError(c *gin.Context, details []string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Code:    CodeValidationError,
		Message: "Validation failed",
		Details: details,
	})
}
