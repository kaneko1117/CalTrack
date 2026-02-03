package repository

import (
	"context"
	"time"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
)

// DailyCalories は日別カロリー集計結果
type DailyCalories struct {
	Date     vo.EatenAt
	Calories vo.Calories
}

// RecordRepository はカロリー記録の永続化を担当するリポジトリインターフェース
type RecordRepository interface {
	// Save はRecordを保存する
	Save(ctx context.Context, record *entity.Record) error
	// FindByUserIDAndDateRange は指定ユーザーの指定日付範囲内のRecordを取得する
	// startTime以上、endTime未満のeatenAtを持つRecordを返す
	// Recordには関連するRecordItemsも含まれる
	FindByUserIDAndDateRange(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error)
	// GetDailyCalories は日別カロリーを取得（グラフ用）
	GetDailyCalories(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) ([]DailyCalories, error)
}
