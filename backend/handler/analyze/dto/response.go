package dto

import (
	"caltrack/usecase"
)

// AnalyzedItemResponse は解析された食品1件のレスポンスDTO
type AnalyzedItemResponse struct {
	Name     string `json:"name"`     // 食品名
	Calories int    `json:"calories"` // カロリー
}

// AnalyzeImageResponse は画像解析レスポンスDTO
type AnalyzeImageResponse struct {
	Items []AnalyzedItemResponse `json:"items"` // 解析された食品リスト
}

// NewAnalyzeImageResponse はUsecaseの出力からレスポンスDTOを生成する
func NewAnalyzeImageResponse(output *usecase.AnalyzeOutput) AnalyzeImageResponse {
	items := make([]AnalyzedItemResponse, len(output.Items))
	for i, item := range output.Items {
		items[i] = AnalyzedItemResponse{
			Name:     item.Name.String(),
			Calories: item.Calories.Value(),
		}
	}

	return AnalyzeImageResponse{
		Items: items,
	}
}
