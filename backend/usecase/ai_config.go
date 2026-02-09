package usecase

// AIConfig はAIサービスの設定を提供するインターフェース
// Usecase層がconfigパッケージに直接依存することを防ぐ
type AIConfig interface {
	// GeminiModelName は使用するGeminiモデル名を返す
	GeminiModelName() string
}
