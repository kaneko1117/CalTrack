package common

import "caltrack/pkg/logger"

// LogError はHandler層のエラーをログ出力する
func LogError(operation string, err error, fields ...any) {
	args := append([]any{"layer", "handler", "operation", operation, "error", err.Error()}, fields...)
	logger.Error("handler error", args...)
}
