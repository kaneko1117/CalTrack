package dto

import (
	"time"

	"caltrack/domain/entity"
)

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
