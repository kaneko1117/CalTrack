package usecase

import (
	"context"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/usecase/service"
)

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

	// 画像解析サービスを呼び出し
	analyzedItems, err := u.imageAnalyzer.Analyze(ctx, imageData, mimeType)
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
