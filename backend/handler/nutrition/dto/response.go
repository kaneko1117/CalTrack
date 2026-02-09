package dto

import (
	"math"

	"caltrack/domain/vo"
	"caltrack/usecase"
)

// PfcResponse はPFC栄養素のレスポンス構造体
type PfcResponse struct {
	Protein       int `json:"protein"`        // タンパク質(g)
	Fat           int `json:"fat"`            // 脂質(g)
	Carbohydrates int `json:"carbohydrates"`  // 炭水化物(g)
	TotalCalories int `json:"total_calories"` // 総カロリー(kcal)
}

// TodayPfcResponse は今日のPFC摂取量と目標PFCのレスポンスDTO
type TodayPfcResponse struct {
	Date       string      `json:"date"`        // 対象日付 (YYYY-MM-DD)
	CurrentPfc PfcResponse `json:"current_pfc"` // 今日のPFC摂取量合計
	TargetPfc  PfcResponse `json:"target_pfc"`  // 目標PFC
}

// pfcToResponse はPfc VOをPfcResponseに変換する
func pfcToResponse(pfc vo.Pfc) PfcResponse {
	// カロリー計算: P * 4 + F * 9 + C * 4
	totalCal := pfc.Protein()*vo.ProteinCalPerGram +
		pfc.Fat()*vo.FatCalPerGram +
		pfc.Carbs()*vo.CarbsCalPerGram

	return PfcResponse{
		Protein:       int(math.Round(pfc.Protein())),
		Fat:           int(math.Round(pfc.Fat())),
		Carbohydrates: int(math.Round(pfc.Carbs())),
		TotalCalories: int(math.Round(totalCal)),
	}
}

// NewTodayPfcResponse はUsecaseの出力からレスポンスDTOを生成する
func NewTodayPfcResponse(output *usecase.TodayPfcOutput) TodayPfcResponse {
	return TodayPfcResponse{
		Date:       output.Date.Format("2006-01-02"),
		CurrentPfc: pfcToResponse(output.CurrentPfc),
		TargetPfc:  pfcToResponse(output.TargetPfc),
	}
}
