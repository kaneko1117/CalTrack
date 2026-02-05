package dto

import (
	"caltrack/usecase/service"
)

// AdviceResponse は栄養アドバイスレスポンスDTO
type AdviceResponse struct {
	Advice string `json:"advice"`
}

// NewAdviceResponse はUsecaseの出力からレスポンスDTOを生成する
func NewAdviceResponse(output *service.NutritionAdviceOutput) AdviceResponse {
	return AdviceResponse{
		Advice: output.Advice,
	}
}
