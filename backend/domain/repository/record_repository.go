package repository

import (
	"context"
	"time"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
)

// RecordRepository はカロリー記録の永続化を担当するリポジトリインターフェース
type RecordRepository interface {
	// Save はRecordを保存する
	Save(ctx context.Context, record *entity.Record) error
	// FindByUserIDAndDateRange は指定ユーザーの指定日付範囲内のRecordを取得する
	// startTime以上、endTime未満のeatenAtを持つRecordを返す
	// Recordには関連するRecordItemsも含まれる
	FindByUserIDAndDateRange(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error)
}
