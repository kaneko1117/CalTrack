package gorm

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
	"caltrack/infrastructure/persistence/gorm/model"
)

type GormAdviceCacheRepository struct {
	db *gorm.DB
}

func NewGormAdviceCacheRepository(db *gorm.DB) *GormAdviceCacheRepository {
	return &GormAdviceCacheRepository{db: db}
}

func (r *GormAdviceCacheRepository) Save(ctx context.Context, cache *entity.AdviceCache) error {
	tx := GetTx(ctx, r.db)
	cacheModel := toAdviceCacheModel(cache)
	if err := tx.Create(&cacheModel).Error; err != nil {
		logError("Save", err, "advice_cache_id", cache.ID().String())
		return err
	}
	return nil
}

func (r *GormAdviceCacheRepository) FindByUserIDAndDate(ctx context.Context, userID vo.UserID, date time.Time) (*entity.AdviceCache, error) {
	tx := GetTx(ctx, r.db)
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var m model.AdviceCache
	err := tx.Where("user_id = ? AND cache_date = ?", userID.String(), normalizedDate).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logError("FindByUserIDAndDate", err, "user_id", userID.String())
		return nil, err
	}
	return toAdviceCacheEntity(&m), nil
}

func (r *GormAdviceCacheRepository) DeleteByUserIDAndDate(ctx context.Context, userID vo.UserID, date time.Time) error {
	tx := GetTx(ctx, r.db)
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	if err := tx.Where("user_id = ? AND cache_date = ?", userID.String(), normalizedDate).Delete(&model.AdviceCache{}).Error; err != nil {
		logError("DeleteByUserIDAndDate", err, "user_id", userID.String())
		return err
	}
	return nil
}

func toAdviceCacheModel(cache *entity.AdviceCache) model.AdviceCache {
	return model.AdviceCache{
		ID:        cache.ID().String(),
		UserID:    cache.UserID().String(),
		CacheDate: cache.CacheDate(),
		Advice:    cache.Advice(),
		CreatedAt: cache.CreatedAt(),
	}
}

func toAdviceCacheEntity(m *model.AdviceCache) *entity.AdviceCache {
	return entity.ReconstructAdviceCache(m.ID, m.UserID, m.CacheDate, m.Advice, m.CreatedAt)
}
