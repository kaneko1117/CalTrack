package service

import (
	"context"

	"caltrack/domain/vo"
)

// AnalyzedItem は画像解析で認識された食品1件を表す
type AnalyzedItem struct {
	Name     vo.ItemName // 食品名
	Calories vo.Calories // カロリー
}

// ImageAnalyzerLogConfig は画像解析のログ出力設定を保持する
type ImageAnalyzerLogConfig struct {
	EnableRequestLog  bool // リクエスト情報のログを有効にする
	EnableResponseLog bool // レスポンス情報のログを有効にする
	EnableTokenLog    bool // トークン使用量のログを有効にする
}

// ImageAnalyzerConfig は画像解析の設定を保持する
type ImageAnalyzerConfig struct {
	ModelName string                 // 使用するAIモデル名
	Prompt    string                 // 解析に使用するプロンプト
	Log       ImageAnalyzerLogConfig // ログ出力設定
}

// ImageAnalyzer は画像から食品を解析するサービスのインターフェース
// Infrastructure層で具体的な実装（Gemini API等）を提供する
type ImageAnalyzer interface {
	// Analyze は画像データを解析し、認識した食品のリストを返す
	Analyze(ctx context.Context, config ImageAnalyzerConfig, imageData string, mimeType string) ([]AnalyzedItem, error)
}
