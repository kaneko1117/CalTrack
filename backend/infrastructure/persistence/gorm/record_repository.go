package gorm

import (
	"context"
	"time"

	"gorm.io/gorm"

	"caltrack/domain/entity"
	"caltrack/domain/repository"
	"caltrack/domain/vo"
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

// FindByUserIDAndDateRange は指定ユーザーの指定日付範囲内のRecordを取得する
// startTime以上、endTime未満のeatenAtを持つRecordを返す
// Recordには関連するRecordItemsも含まれる
func (r *GormRecordRepository) FindByUserIDAndDateRange(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
	tx := GetTx(ctx, r.db)
	var models []model.Record
	err := tx.Where("user_id = ? AND eaten_at >= ? AND eaten_at < ?", userID.String(), startTime, endTime).
		Preload("Items").
		Order("eaten_at ASC").
		Find(&models).Error
	if err != nil {
		logError("FindByUserIDAndDateRange", err, "user_id", userID.String())
		return nil, err
	}

	records := make([]*entity.Record, len(models))
	for i, m := range models {
		records[i] = toRecordEntity(&m)
	}
	return records, nil
}

// GetDailyCalories は日別カロリーを取得（グラフ用）
func (r *GormRecordRepository) GetDailyCalories(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
	tx := GetTx(ctx, r.db)

	// 期間の計算（今日を含めて過去N日間）
	now := time.Now()
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
	startOfPeriod := endOfDay.AddDate(0, 0, -period.Days())

	// 日別カロリー集計クエリ
	type dailySum struct {
		Date          time.Time
		TotalCalories int
	}
	var results []dailySum

	// records と record_items を JOIN して日別集計
	err := tx.Table("records").
		Select("DATE(eaten_at) as date, COALESCE(SUM(record_items.calories), 0) as total_calories").
		Joins("LEFT JOIN record_items ON records.id = record_items.record_id").
		Where("records.user_id = ? AND records.eaten_at >= ? AND records.eaten_at < ?", userID.String(), startOfPeriod, endOfDay).
		Group("DATE(eaten_at)").
		Order("date ASC").
		Find(&results).Error
	if err != nil {
		logError("GetDailyCalories", err, "user_id", userID.String())
		return nil, err
	}

	// 結果をVO形式に変換
	dailyCalories := make([]repository.DailyCalories, len(results))
	for i, r := range results {
		eatenAt := vo.ReconstructEatenAt(r.Date)
		calories := vo.ReconstructCalories(r.TotalCalories)
		dailyCalories[i] = repository.DailyCalories{
			Date:     eatenAt,
			Calories: calories,
		}
	}

	return dailyCalories, nil
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

// toRecordEntity はGORMモデルをエンティティに変換する
func toRecordEntity(m *model.Record) *entity.Record {
	items := make([]entity.RecordItem, len(m.Items))
	for i, itemModel := range m.Items {
		items[i] = *entity.ReconstructRecordItem(
			itemModel.ID,
			itemModel.RecordID,
			itemModel.Name,
			itemModel.Calories,
		)
	}
	return entity.ReconstructRecord(
		m.ID,
		m.UserID,
		m.EatenAt,
		m.CreatedAt,
		items,
	)
}
