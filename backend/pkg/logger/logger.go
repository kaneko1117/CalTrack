package logger

import (
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/lmittmann/tint"
)

var (
	defaultLogger *slog.Logger
	once          sync.Once
)

// Init はロガーを初期化する
// 環境変数ENVに応じて出力形式を切り替える
// - production: JSON形式（本番向け）
// - それ以外: カラー形式（開発向け）
func Init() {
	once.Do(func() {
		initLogger()
	})
}

// initLogger は内部でロガーを初期化する
func initLogger() {
	var handler slog.Handler

	if os.Getenv("ENV") == "production" {
		// 本番環境: JSON出力
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	} else {
		// 開発環境: カラー出力
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.DateTime,
			AddSource:  false,
		})
	}

	defaultLogger = slog.New(handler)
	slog.SetDefault(defaultLogger)
}

// getLogger はデフォルトロガーを取得する
// 初期化されていない場合は自動的に初期化する
func getLogger() *slog.Logger {
	if defaultLogger == nil {
		Init()
	}
	return defaultLogger
}

// Error はエラーレベルのログを出力する
func Error(msg string, args ...any) {
	getLogger().Error(msg, args...)
}

// Warn は警告レベルのログを出力する
func Warn(msg string, args ...any) {
	getLogger().Warn(msg, args...)
}

// Info は情報レベルのログを出力する
func Info(msg string, args ...any) {
	getLogger().Info(msg, args...)
}

// Debug はデバッグレベルのログを出力する
func Debug(msg string, args ...any) {
	getLogger().Debug(msg, args...)
}

// With は追加のフィールドを持つロガーを返す
func With(args ...any) *slog.Logger {
	return getLogger().With(args...)
}
