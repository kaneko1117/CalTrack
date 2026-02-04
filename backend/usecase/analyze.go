package usecase

import (
	"context"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/usecase/service"
)

// 画像解析に使用するAIモデル名
const geminiModelName = "gemini-2.0-flash"

// 画像解析に使用するプロンプト
const analyzePrompt = `この画像に写っている食品を分析し、以下のJSON形式で回答してください。
食品が複数ある場合は全て列挙してください。
カロリーは一般的な1人前の量で推定してください。

回答形式（JSON配列のみを返してください。マークダウンのコードブロックは使わないでください）:
[
  {"name": "食品名", "calories": カロリー数値}
]

例:
[
  {"name": "白ご飯", "calories": 235},
  {"name": "味噌汁", "calories": 40}
]

画像に食品が写っていない場合は空の配列を返してください: []`

// AnalyzeOutput は画像解析結果の出力構造体
type AnalyzeOutput struct {
	Items []AnalyzedItemOutput // 解析された食品リスト
}

// AnalyzedItemOutput は解析された食品1件の出力
type AnalyzedItemOutput struct {
	Name     vo.ItemName // 食品名
	Calories vo.Calories // カロリー
}

// AnalyzeUsecase は画像解析に関するユースケースを提供する
type AnalyzeUsecase struct {
	imageAnalyzer service.ImageAnalyzer
}

// NewAnalyzeUsecase は AnalyzeUsecase のインスタンスを生成する
func NewAnalyzeUsecase(imageAnalyzer service.ImageAnalyzer) *AnalyzeUsecase {
	return &AnalyzeUsecase{
		imageAnalyzer: imageAnalyzer,
	}
}

// AnalyzeImage は画像から食品を解析し、カロリー情報を返す
func (u *AnalyzeUsecase) AnalyzeImage(ctx context.Context, imageData string, mimeType string) (*AnalyzeOutput, error) {
	// 入力バリデーション
	if imageData == "" {
		return nil, domainErrors.ErrImageDataRequired
	}
	if mimeType == "" {
		return nil, domainErrors.ErrMimeTypeRequired
	}

	// 解析設定を構築（ビジネスロジックとしてUsecase層で管理）
	config := service.ImageAnalyzerConfig{
		ModelName: geminiModelName,
		Prompt:    analyzePrompt,
	}

	// 画像解析サービスを呼び出し
	analyzedItems, err := u.imageAnalyzer.Analyze(ctx, config, imageData, mimeType)
	if err != nil {
		logError("AnalyzeImage", err, "mimeType", mimeType)
		return nil, err
	}

	// 結果が空の場合
	if len(analyzedItems) == 0 {
		return nil, domainErrors.ErrNoFoodDetected
	}

	// 出力に変換
	outputItems := make([]AnalyzedItemOutput, len(analyzedItems))
	for i, item := range analyzedItems {
		outputItems[i] = AnalyzedItemOutput{
			Name:     item.Name,
			Calories: item.Calories,
		}
	}

	return &AnalyzeOutput{
		Items: outputItems,
	}, nil
}
