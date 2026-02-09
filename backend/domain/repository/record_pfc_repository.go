package repository

import (
	"context"
	"time"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
)

type RecordPfcRepository interface {
	Save(ctx context.Context, recordPfc *entity.RecordPfc) error
	FindByRecordID(ctx context.Context, recordID vo.RecordID) (*entity.RecordPfc, error)
	FindByRecordIDs(ctx context.Context, recordIDs []vo.RecordID) ([]*entity.RecordPfc, error)
	GetDailyPfc(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) (vo.DailyPfc, error)
}
