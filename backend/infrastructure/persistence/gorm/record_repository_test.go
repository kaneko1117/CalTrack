package gorm_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	gormPkg "caltrack/infrastructure/persistence/gorm"
)

// ============================================================================
// Save テスト
// ============================================================================

func TestGormRecordRepository_Save(t *testing.T) {
	t.Run("正常系_RecordとItemsが保存される", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordRepository(db)
		ctx := context.Background()

		user := testUser(t)
		eatenAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		record := testRecordWithItem(t, user.ID(), eatenAt, "ランチ", 500)

		// GORMのCreate()はRecordとRecordItemsを一括保存
		// 1. Recordのインサート
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `records`")).
			WithArgs(
				record.ID().String(),
				record.UserID().String(),
				record.EatenAt().Time(),
				sqlmock.AnyArg(), // created_at
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		// 2. RecordItemsのインサート（Itemsが含まれている場合）
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_items`")).
			WithArgs(
				record.Items()[0].ID().String(),
				record.Items()[0].RecordID().String(),
				record.Items()[0].Name().String(),
				record.Items()[0].Calories().Value(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Save実行
		err := repo.Save(ctx, record)
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	})

	t.Run("異常系_DBエラーで保存失敗", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordRepository(db)
		ctx := context.Background()

		user := testUser(t)
		eatenAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		record := testRecord(t, user.ID(), eatenAt)

		// INSERT失敗をシミュレート
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `records`")).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		// Save実行
		err := repo.Save(ctx, record)
		if err == nil {
			t.Error("Save() should fail with db error")
		}
	})
}

// ============================================================================
// FindByUserIDAndDateRange テスト
// ============================================================================

func TestGormRecordRepository_FindByUserIDAndDateRange(t *testing.T) {
	t.Run("正常系_日付範囲内のRecordが取得できる", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordRepository(db)
		ctx := context.Background()

		user := testUser(t)
		startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
		eatenAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

		record := testRecordWithItem(t, user.ID(), eatenAt, "ランチ", 500)

		// 1. メインSELECT（WHERE user_id = ? AND eaten_at >= ? AND eaten_at < ? ORDER BY eaten_at ASC）
		rows := sqlmock.NewRows(recordColumns()).
			AddRow(
				record.ID().String(),
				record.UserID().String(),
				record.EatenAt().Time(),
				record.CreatedAt(),
			)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `records` WHERE user_id = ? AND eaten_at >= ? AND eaten_at < ?")).
			WithArgs(user.ID().String(), startTime, endTime).
			WillReturnRows(rows)

		// 2. Preload用SELECT（record_items WHERE record_id IN (?)）
		itemRows := sqlmock.NewRows(recordItemColumns()).
			AddRow(
				record.Items()[0].ID().String(),
				record.Items()[0].RecordID().String(),
				record.Items()[0].Name().String(),
				record.Items()[0].Calories().Value(),
			)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_items` WHERE `record_items`.`record_id` = ?")).
			WithArgs(record.ID().String()).
			WillReturnRows(itemRows)

		// FindByUserIDAndDateRange実行
		found, err := repo.FindByUserIDAndDateRange(ctx, user.ID(), startTime, endTime)
		if err != nil {
			t.Fatalf("FindByUserIDAndDateRange() error = %v", err)
		}
		if len(found) != 1 {
			t.Fatalf("expected 1 record, got %d", len(found))
		}
		if !found[0].ID().Equals(record.ID()) {
			t.Errorf("ID = %v, want %v", found[0].ID(), record.ID())
		}
		if len(found[0].Items()) != 1 {
			t.Errorf("expected 1 item, got %d", len(found[0].Items()))
		}
	})

	t.Run("正常系_該当なしで空配列が返る", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordRepository(db)
		ctx := context.Background()

		user := testUser(t)
		startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

		// 空のrowsを返す
		rows := sqlmock.NewRows(recordColumns())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `records` WHERE user_id = ? AND eaten_at >= ? AND eaten_at < ?")).
			WithArgs(user.ID().String(), startTime, endTime).
			WillReturnRows(rows)

		// Preload用クエリは発行されない（メインクエリが空なので）

		// FindByUserIDAndDateRange実行
		found, err := repo.FindByUserIDAndDateRange(ctx, user.ID(), startTime, endTime)
		if err != nil {
			t.Fatalf("FindByUserIDAndDateRange() error = %v", err)
		}
		if len(found) != 0 {
			t.Errorf("expected empty array, got %d records", len(found))
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordRepository(db)
		ctx := context.Background()

		user := testUser(t)
		startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `records` WHERE user_id = ? AND eaten_at >= ? AND eaten_at < ?")).
			WithArgs(user.ID().String(), startTime, endTime).
			WillReturnError(errors.New("db error"))

		// FindByUserIDAndDateRange実行
		found, err := repo.FindByUserIDAndDateRange(ctx, user.ID(), startTime, endTime)
		if err == nil {
			t.Error("FindByUserIDAndDateRange() should fail with db error")
		}
		if found != nil {
			t.Error("found records should be nil on error")
		}
	})
}

// ============================================================================
// GetDailyCalories テスト
// ============================================================================

func TestGormRecordRepository_GetDailyCalories(t *testing.T) {
	t.Run("正常系_日別カロリー集計が取得できる", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordRepository(db)
		ctx := context.Background()

		user := testUser(t)
		period := testStatisticsPeriod(t, "week") // 7日間

		// 期間の計算（record_repository.goの実装と同じ）
		now := time.Now()
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
		startOfPeriod := endOfDay.AddDate(0, 0, -period.Days())

		// 集計結果: DATE(eaten_at), COALESCE(SUM(record_items.calories), 0)
		rows := sqlmock.NewRows([]string{"date", "total_calories"}).
			AddRow(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), 1500).
			AddRow(time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), 2000)

		// JOIN + GROUP BY + SUM のクエリ
		mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE(eaten_at) as date, COALESCE(SUM(record_items.calories), 0) as total_calories FROM `records` LEFT JOIN record_items ON records.id = record_items.record_id WHERE")).
			WithArgs(user.ID().String(), startOfPeriod, endOfDay).
			WillReturnRows(rows)

		// GetDailyCalories実行
		result, err := repo.GetDailyCalories(ctx, user.ID(), period)
		if err != nil {
			t.Fatalf("GetDailyCalories() error = %v", err)
		}
		if len(result) != 2 {
			t.Fatalf("expected 2 daily records, got %d", len(result))
		}
		if result[0].Calories.Value() != 1500 {
			t.Errorf("expected calories 1500, got %d", result[0].Calories.Value())
		}
		if result[1].Calories.Value() != 2000 {
			t.Errorf("expected calories 2000, got %d", result[1].Calories.Value())
		}
	})

	t.Run("正常系_該当なしで空配列が返る", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordRepository(db)
		ctx := context.Background()

		user := testUser(t)
		period := testStatisticsPeriod(t, "week")

		now := time.Now()
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
		startOfPeriod := endOfDay.AddDate(0, 0, -period.Days())

		// 空の集計結果
		rows := sqlmock.NewRows([]string{"date", "total_calories"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE(eaten_at) as date, COALESCE(SUM(record_items.calories), 0) as total_calories FROM `records` LEFT JOIN record_items ON records.id = record_items.record_id WHERE")).
			WithArgs(user.ID().String(), startOfPeriod, endOfDay).
			WillReturnRows(rows)

		result, err := repo.GetDailyCalories(ctx, user.ID(), period)
		if err != nil {
			t.Fatalf("GetDailyCalories() error = %v", err)
		}
		if len(result) != 0 {
			t.Errorf("expected empty array, got %d records", len(result))
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordRepository(db)
		ctx := context.Background()

		user := testUser(t)
		period := testStatisticsPeriod(t, "week")

		now := time.Now()
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
		startOfPeriod := endOfDay.AddDate(0, 0, -period.Days())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT DATE(eaten_at) as date, COALESCE(SUM(record_items.calories), 0) as total_calories FROM `records` LEFT JOIN record_items ON records.id = record_items.record_id WHERE")).
			WithArgs(user.ID().String(), startOfPeriod, endOfDay).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetDailyCalories(ctx, user.ID(), period)
		if err == nil {
			t.Error("GetDailyCalories() should fail with db error")
		}
		if result != nil {
			t.Error("result should be nil on error")
		}
	})
}
