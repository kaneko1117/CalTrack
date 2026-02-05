package repository

import (
	"context"
	"time"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
)

// AdviceCacheRepository はアドバイスキャッシュの永続化を担当するリポジトリインターフェース
type AdviceCacheRepository interface {
	Save(ctx context.Context, cache *entity.AdviceCache) error
	FindByUserIDAndDate(ctx context.Context, userID vo.UserID, date time.Time) (*entity.AdviceCache, error)
	DeleteByUserIDAndDate(ctx context.Context, userID vo.UserID, date time.Time) error
}
