package logger

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Middleware はリクエストログを出力するGinミドルウェア
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// リクエスト処理
		c.Next()

		// レスポンス後の処理
		latency := time.Since(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()

		// ステータスコードに応じてログレベルを変更
		if status >= 500 {
			Error("request failed",
				"status", status,
				"method", method,
				"path", path,
				"latency_ms", latency.Milliseconds(),
				"client_ip", clientIP,
			)
		} else if status >= 400 {
			Warn("client error",
				"status", status,
				"method", method,
				"path", path,
				"latency_ms", latency.Milliseconds(),
				"client_ip", clientIP,
			)
		} else {
			Info("request completed",
				"status", status,
				"method", method,
				"path", path,
				"latency_ms", latency.Milliseconds(),
				"client_ip", clientIP,
			)
		}
	}
}
