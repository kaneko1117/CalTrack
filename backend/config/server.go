package config

import (
	"os"
	"strings"
)

const defaultCORSAllowOrigins = "http://localhost:5173"

// GetCORSAllowOrigins は CORS で許可するオリジンを取得する
// 環境変数 CORS_ALLOW_ORIGINS から取得し、未設定の場合はデフォルト値を使用
// 複数オリジンはカンマ区切りで指定
func GetCORSAllowOrigins() []string {
	originsStr := os.Getenv("CORS_ALLOW_ORIGINS")
	if originsStr == "" {
		originsStr = defaultCORSAllowOrigins
	}

	origins := strings.Split(originsStr, ",")
	result := make([]string, 0, len(origins))
	for _, origin := range origins {
		trimmed := strings.TrimSpace(origin)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
