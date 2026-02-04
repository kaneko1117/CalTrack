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

// ImageAnalyzer は画像から食品を解析するサービスのインターフェース
// Infrastructure層で具体的な実装（Gemini API等）を提供する
type ImageAnalyzer interface {
	// Analyze は画像データを解析し、認識した食品のリストを返す
	// imageData: Base64エンコードされた画像データ
	// mimeType: 画像のMIMEタイプ（例: "image/jpeg", "image/png"）
	Analyze(ctx context.Context, imageData string, mimeType string) ([]AnalyzedItem, error)
}
