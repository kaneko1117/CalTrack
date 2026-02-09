package dto

import (
	"time"

	"caltrack/usecase"
	"caltrack/usecase/service"
)

// AdviceResponse は栄養アドバイスのレスポンスDTO
type AdviceResponse struct {
	Advice string `json:"advice"` // アドバイス内容
}

// NewAdviceResponse はAdviceOutputからレスポンスDTOを生成する
func NewAdviceResponse(output *service.NutritionAdviceOutput) AdviceResponse {
	return AdviceResponse{
		Advice: output.Advice,
	}
}

// TodayPfcResponse は今日1日のPFC摂取量と目標のレスポンスDTO
type TodayPfcResponse struct {
	Date    time.Time `json:"date"`
	Current PfcDTO    `json:"current"`
	Target  PfcDTO    `json:"target"`
}

// PfcDTO はPFC値のDTO
type PfcDTO struct {
	Protein float64 `json:"protein"`
	Fat     float64 `json:"fat"`
	Carbs   float64 `json:"carbs"`
}

// NewTodayPfcResponse はUsecaseの出力からレスポンスDTOを生成する
func NewTodayPfcResponse(output *usecase.TodayPfcOutput) TodayPfcResponse {
	return TodayPfcResponse{
		Date: output.Date,
		Current: PfcDTO{
			Protein: output.CurrentPfc.Protein(),
			Fat:     output.CurrentPfc.Fat(),
			Carbs:   output.CurrentPfc.Carbs(),
		},
		Target: PfcDTO{
			Protein: output.TargetPfc.Protein(),
			Fat:     output.TargetPfc.Fat(),
			Carbs:   output.TargetPfc.Carbs(),
		},
	}
}
