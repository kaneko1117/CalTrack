package entity

import (
	"time"

	"caltrack/domain/vo"
)

// AdviceCache はAIアドバイスのキャッシュを表すEntity
type AdviceCache struct {
	id        vo.AdviceCacheID
	userID    vo.UserID
	cacheDate time.Time
	advice    string
	createdAt time.Time
}

// NewAdviceCache は新しいAdviceCacheを生成する
func NewAdviceCache(userID vo.UserID, cacheDate time.Time, advice string) *AdviceCache {
	normalizedDate := time.Date(cacheDate.Year(), cacheDate.Month(), cacheDate.Day(), 0, 0, 0, 0, cacheDate.Location())

	return &AdviceCache{
		id:        vo.NewAdviceCacheID(),
		userID:    userID,
		cacheDate: normalizedDate,
		advice:    advice,
		createdAt: time.Now(),
	}
}

// ReconstructAdviceCache はDBからAdviceCacheを復元する
func ReconstructAdviceCache(idStr, userIDStr string, cacheDate time.Time, advice string, createdAt time.Time) *AdviceCache {
	return &AdviceCache{
		id:        vo.ReconstructAdviceCacheID(idStr),
		userID:    vo.ReconstructUserID(userIDStr),
		cacheDate: cacheDate,
		advice:    advice,
		createdAt: createdAt,
	}
}

func (a *AdviceCache) ID() vo.AdviceCacheID { return a.id }
func (a *AdviceCache) UserID() vo.UserID    { return a.userID }
func (a *AdviceCache) CacheDate() time.Time { return a.cacheDate }
func (a *AdviceCache) Advice() string       { return a.advice }
func (a *AdviceCache) CreatedAt() time.Time { return a.createdAt }
