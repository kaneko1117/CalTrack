package gorm

import (
	"context"

	"gorm.io/gorm"

	"caltrack/domain/entity"
	"caltrack/infrastructure/persistence/gorm/model"
)

// GormRecordRepository はRecordRepositoryのGORM実装
type GormRecordRepository struct {
	db *gorm.DB
}

// NewGormRecordRepository は新しいGormRecordRepositoryを生成する
func NewGormRecordRepository(db *gorm.DB) *GormRecordRepository {
	return &GormRecordRepository{db: db}
}

// Save はRecordを保存する
func (r *GormRecordRepository) Save(ctx context.Context, record *entity.Record) error {
	tx := GetTx(ctx, r.db)

	recordModel := toRecordModel(record)
	if err := tx.Create(&recordModel).Error; err != nil {
		logError("Save", err, "record_id", record.ID().String())
		return err
	}

	return nil
}

// toRecordModel はエンティティをGORMモデルに変換する
func toRecordModel(record *entity.Record) model.Record {
	items := record.Items()
	itemModels := make([]model.RecordItem, len(items))
	for i, item := range items {
		itemModels[i] = model.RecordItem{
			ID:       item.ID().String(),
			RecordID: item.RecordID().String(),
			Name:     item.Name().String(),
			Calories: item.Calories().Value(),
		}
	}

	return model.Record{
		ID:        record.ID().String(),
		UserID:    record.UserID().String(),
		EatenAt:   record.EatenAt().Time(),
		CreatedAt: record.CreatedAt(),
		Items:     itemModels,
	}
}
