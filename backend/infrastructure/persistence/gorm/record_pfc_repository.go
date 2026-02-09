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

// GormRecordPfcRepository はRecordPfcRepositoryのGORM実装
type GormRecordPfcRepository struct {
	db *gorm.DB
}

// NewGormRecordPfcRepository は新しいGormRecordPfcRepositoryを生成する
func NewGormRecordPfcRepository(db *gorm.DB) *GormRecordPfcRepository {
	return &GormRecordPfcRepository{db: db}
}

// Save はRecordPfcを保存する
func (r *GormRecordPfcRepository) Save(ctx context.Context, recordPfc *entity.RecordPfc) error {
	tx := GetTx(ctx, r.db)

	recordPfcModel := toRecordPfcModel(recordPfc)
	if err := tx.Create(&recordPfcModel).Error; err != nil {
		logError("Save", err, "record_pfc_id", recordPfc.ID().String())
		return err
	}

	return nil
}

// FindByRecordID は指定RecordIDのRecordPfcを取得する（見つからなければnil, nil）
func (r *GormRecordPfcRepository) FindByRecordID(ctx context.Context, recordID vo.RecordID) (*entity.RecordPfc, error) {
	tx := GetTx(ctx, r.db)

	var m model.RecordPfc
	err := tx.Where("record_id = ?", recordID.String()).First(&m).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 見つからない場合はnilを返す（エラーではない）
			return nil, nil
		}
		logError("FindByRecordID", err, "record_id", recordID.String())
		return nil, err
	}

	return toRecordPfcEntity(&m), nil
}

// FindByRecordIDs は複数のRecordIDに対応するRecordPfcを一括取得する
func (r *GormRecordPfcRepository) FindByRecordIDs(ctx context.Context, recordIDs []vo.RecordID) ([]*entity.RecordPfc, error) {
	if len(recordIDs) == 0 {
		return []*entity.RecordPfc{}, nil
	}

	tx := GetTx(ctx, r.db)

	// RecordIDを文字列のスライスに変換
	recordIDStrs := make([]string, len(recordIDs))
	for i, id := range recordIDs {
		recordIDStrs[i] = id.String()
	}

	var models []model.RecordPfc
	err := tx.Where("record_id IN ?", recordIDStrs).Find(&models).Error
	if err != nil {
		logError("FindByRecordIDs", err, "count", len(recordIDs))
		return nil, err
	}

	recordPfcs := make([]*entity.RecordPfc, len(models))
	for i, m := range models {
		recordPfcs[i] = toRecordPfcEntity(&m)
	}

	return recordPfcs, nil
}

// GetDailyPfc は指定日時範囲のPFC合計を取得する
func (r *GormRecordPfcRepository) GetDailyPfc(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) (vo.DailyPfc, error) {
	tx := GetTx(ctx, r.db)

	type pfcSum struct {
		TotalProtein float64
		TotalFat     float64
		TotalCarbs   float64
	}
	var result pfcSum

	err := tx.Table("records").
		Select("COALESCE(SUM(record_pfcs.protein), 0) as total_protein, COALESCE(SUM(record_pfcs.fat), 0) as total_fat, COALESCE(SUM(record_pfcs.carbs), 0) as total_carbs").
		Joins("INNER JOIN record_pfcs ON records.id = record_pfcs.record_id").
		Where("records.user_id = ? AND records.eaten_at >= ? AND records.eaten_at < ?", userID.String(), startTime, endTime).
		Scan(&result).Error
	if err != nil {
		logError("GetDailyPfc", err, "user_id", userID.String())
		return vo.DailyPfc{}, err
	}

	return vo.NewDailyPfc(result.TotalProtein, result.TotalFat, result.TotalCarbs), nil
}

// toRecordPfcModel はエンティティをGORMモデルに変換する
func toRecordPfcModel(recordPfc *entity.RecordPfc) model.RecordPfc {
	return model.RecordPfc{
		ID:       recordPfc.ID().String(),
		RecordID: recordPfc.RecordID().String(),
		Protein:  recordPfc.Protein(),
		Fat:      recordPfc.Fat(),
		Carbs:    recordPfc.Carbs(),
	}
}

// toRecordPfcEntity はGORMモデルをエンティティに変換する
func toRecordPfcEntity(m *model.RecordPfc) *entity.RecordPfc {
	return entity.ReconstructRecordPfc(
		m.ID,
		m.RecordID,
		m.Protein,
		m.Fat,
		m.Carbs,
	)
}
