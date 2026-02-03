package dto

import (
	"time"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
)

// CreateRecordRequest はカロリー記録作成リクエストDTO
type CreateRecordRequest struct {
	EatenAt string              `json:"eatenAt"`
	Items   []RecordItemRequest `json:"items"`
}

// RecordItemRequest は記録明細リクエストDTO
type RecordItemRequest struct {
	Name     string `json:"name"`
	Calories int    `json:"calories"`
}

// ToDomain はリクエストをEntityに変換する
func (r CreateRecordRequest) ToDomain(userIDStr string) (*entity.Record, error, []error) {
	// UserIDの復元
	userID := vo.ReconstructUserID(userIDStr)

	// 日時のパース
	eatenAtTime, parseErr := time.Parse(time.RFC3339, r.EatenAt)
	if parseErr != nil {
		return nil, parseErr, nil
	}

	// Record作成
	record, err := entity.NewRecord(userID, eatenAtTime)
	if err != nil {
		return nil, nil, []error{err}
	}

	// Items追加
	var validationErrs []error
	for _, item := range r.Items {
		if err := record.AddItem(item.Name, item.Calories); err != nil {
			validationErrs = append(validationErrs, err)
		}
	}

	if len(validationErrs) > 0 {
		return nil, nil, validationErrs
	}

	return record, nil, nil
}

// GetStatisticsRequest は統計データ取得リクエストDTO
type GetStatisticsRequest struct {
	Period string `form:"period"` // クエリパラメータ: week または month
}

// ToDomain はリクエストを StatisticsPeriod VOに変換する
func (r GetStatisticsRequest) ToDomain() (vo.StatisticsPeriod, error) {
	return vo.NewStatisticsPeriod(r.Period)
}
