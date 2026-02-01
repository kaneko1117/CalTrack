package repository

import (
	"context"

	"caltrack/domain/entity"
)

// RecordRepository はカロリー記録の永続化を担当するリポジトリインターフェース
type RecordRepository interface {
	// Save はRecordを保存する
	Save(ctx context.Context, record *entity.Record) error
}
