package usecase

import "caltrack/pkg/logger"

// logError はUsecase層のエラーをログ出力する
func logError(operation string, err error, fields ...any) {
	args := append([]any{"layer", "usecase", "operation", operation, "error", err.Error()}, fields...)
	logger.Error("usecase error", args...)
}

// logWarn はUsecase層の警告をログ出力する
func logWarn(operation string, message string, fields ...any) {
	args := append([]any{"layer", "usecase", "operation", operation}, fields...)
	logger.Warn(message, args...)
}
