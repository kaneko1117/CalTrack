package gorm

import "caltrack/pkg/logger"

// logError はRepository層のエラーをログ出力する
func logError(operation string, err error, fields ...any) {
	args := append([]any{"layer", "repository", "operation", operation, "error", err.Error()}, fields...)
	logger.Error("repository error", args...)
}
