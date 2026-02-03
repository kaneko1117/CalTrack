package dto

import (
	"time"

	"caltrack/domain/entity"
	"caltrack/usecase"
)

// DailyStatisticsResponse は日別統計データレスポンスDTO
type DailyStatisticsResponse struct {
	Date          string `json:"date"`          // YYYY-MM-DD
	TotalCalories int    `json:"totalCalories"` // その日の合計カロリー
}

// StatisticsResponse は統計データレスポンスDTO
type StatisticsResponse struct {
	Period          string                    `json:"period"`          // week/month
	TargetCalories  int                       `json:"targetCalories"`  // 目標カロリー
	AverageCalories int                       `json:"averageCalories"` // 平均カロリー
	TotalDays       int                       `json:"totalDays"`       // 期間の日数
	AchievedDays    int                       `json:"achievedDays"`    // 達成日数
	OverDays        int                       `json:"overDays"`        // 超過日数
	DailyStatistics []DailyStatisticsResponse `json:"dailyStatistics"` // 日別データ
}

// NewStatisticsResponse はUsecaseの出力からレスポンスDTOを生成する
func NewStatisticsResponse(output *usecase.StatisticsOutput) StatisticsResponse {
	dailyStats := make([]DailyStatisticsResponse, len(output.DailyStatistics))
	for i, daily := range output.DailyStatistics {
		dailyStats[i] = DailyStatisticsResponse{
			Date:          daily.Date.Time().Format("2006-01-02"),
			TotalCalories: daily.TotalCalories.Value(),
		}
	}

	return StatisticsResponse{
		Period:          output.Period.String(),
		TargetCalories:  output.TargetCalories.Value(),
		AverageCalories: output.AverageCalories.Value(),
		TotalDays:       output.TotalDays,
		AchievedDays:    output.AchievedDays,
		OverDays:        output.OverDays,
		DailyStatistics: dailyStats,
	}
}

// CreateRecordResponse はカロリー記録作成レスポンスDTO
type CreateRecordResponse struct {
	RecordID      string               `json:"recordId"`
	EatenAt       string               `json:"eatenAt"`
	TotalCalories int                  `json:"totalCalories"`
	Items         []RecordItemResponse `json:"items"`
}

// RecordItemResponse は記録明細レスポンスDTO
type RecordItemResponse struct {
	ItemID   string `json:"itemId"`
	Name     string `json:"name"`
	Calories int    `json:"calories"`
}

// NewCreateRecordResponse はEntityからレスポンスDTOを生成する
func NewCreateRecordResponse(record *entity.Record) CreateRecordResponse {
	items := make([]RecordItemResponse, len(record.Items()))
	for i, item := range record.Items() {
		items[i] = RecordItemResponse{
			ItemID:   item.ID().String(),
			Name:     item.Name().String(),
			Calories: item.Calories().Value(),
		}
	}

	return CreateRecordResponse{
		RecordID:      record.ID().String(),
		EatenAt:       record.EatenAt().Time().Format(time.RFC3339),
		TotalCalories: record.TotalCalories(),
		Items:         items,
	}
}

// TodayCaloriesResponse は今日の摂取カロリーレスポンスDTO
type TodayCaloriesResponse struct {
	Date           string           `json:"date"`
	TotalCalories  int              `json:"totalCalories"`
	TargetCalories int              `json:"targetCalories"`
	Difference     int              `json:"difference"`
	Records        []RecordResponse `json:"records"`
}

// RecordResponse は記録レスポンスDTO
type RecordResponse struct {
	ID      string               `json:"id"`
	EatenAt string               `json:"eatenAt"`
	Items   []RecordItemResponse `json:"items"`
}

// NewTodayCaloriesResponse はUsecaseの出力からレスポンスDTOを生成する
func NewTodayCaloriesResponse(output *usecase.TodayCaloriesOutput) TodayCaloriesResponse {
	records := make([]RecordResponse, len(output.Records))
	for i, record := range output.Records {
		items := make([]RecordItemResponse, len(record.Items()))
		for j, item := range record.Items() {
			items[j] = RecordItemResponse{
				ItemID:   item.ID().String(),
				Name:     item.Name().String(),
				Calories: item.Calories().Value(),
			}
		}
		records[i] = RecordResponse{
			ID:      record.ID().String(),
			EatenAt: record.EatenAt().Time().Format(time.RFC3339),
			Items:   items,
		}
	}

	return TodayCaloriesResponse{
		Date:           output.Date.Format("2006-01-02"),
		TotalCalories:  output.TotalCalories,
		TargetCalories: output.TargetCalories,
		Difference:     output.Difference,
		Records:        records,
	}
}
