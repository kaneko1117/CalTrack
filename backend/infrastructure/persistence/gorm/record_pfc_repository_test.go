package gorm_test

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	gormPkg "caltrack/infrastructure/persistence/gorm"

	"caltrack/domain/vo"
)

// ============================================================================
// Save テスト
// ============================================================================

func TestGormRecordPfcRepository_Save(t *testing.T) {
	t.Run("正常系_RecordPfcが保存される", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordPfcRepository(db)
		ctx := context.Background()

		user := testUser(t)
		eatenAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		record := testRecord(t, user.ID(), eatenAt)
		recordPfc := testRecordPfc(t, record.ID(), 30.0, 20.0, 50.0)

		// GORMのCreate()はINSERTを実行
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_pfcs`")).
			WithArgs(
				recordPfc.ID().String(),
				recordPfc.RecordID().String(),
				recordPfc.Protein(),
				recordPfc.Fat(),
				recordPfc.Carbs(),
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Save実行
		err := repo.Save(ctx, recordPfc)
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	})

	t.Run("異常系_DBエラーで保存失敗", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordPfcRepository(db)
		ctx := context.Background()

		user := testUser(t)
		eatenAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		record := testRecord(t, user.ID(), eatenAt)
		recordPfc := testRecordPfc(t, record.ID(), 30.0, 20.0, 50.0)

		// INSERT失敗をシミュレート
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `record_pfcs`")).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		// Save実行
		err := repo.Save(ctx, recordPfc)
		if err == nil {
			t.Error("Save() should fail with db error")
		}
	})
}

// ============================================================================
// FindByRecordID テスト
// ============================================================================

func TestGormRecordPfcRepository_FindByRecordID(t *testing.T) {
	t.Run("正常系_RecordIDでRecordPfcが見つかる", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordPfcRepository(db)
		ctx := context.Background()

		user := testUser(t)
		eatenAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		record := testRecord(t, user.ID(), eatenAt)
		recordPfc := testRecordPfc(t, record.ID(), 30.0, 20.0, 50.0)

		// GORMのFirst()は ORDER BY id LIMIT 1 を付加
		rows := sqlmock.NewRows(recordPfcColumns()).
			AddRow(
				recordPfc.ID().String(),
				recordPfc.RecordID().String(),
				recordPfc.Protein(),
				recordPfc.Fat(),
				recordPfc.Carbs(),
			)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_pfcs` WHERE record_id = ?")).
			WithArgs(record.ID().String(), 1). // First()は最後に1を追加
			WillReturnRows(rows)

		// FindByRecordID実行
		found, err := repo.FindByRecordID(ctx, record.ID())
		if err != nil {
			t.Fatalf("FindByRecordID() error = %v", err)
		}
		if found == nil {
			t.Fatal("found recordPfc should not be nil")
		}
		if !found.ID().Equals(recordPfc.ID()) {
			t.Errorf("ID = %v, want %v", found.ID(), recordPfc.ID())
		}
		if found.Protein() != recordPfc.Protein() {
			t.Errorf("Protein = %v, want %v", found.Protein(), recordPfc.Protein())
		}
		if found.Fat() != recordPfc.Fat() {
			t.Errorf("Fat = %v, want %v", found.Fat(), recordPfc.Fat())
		}
		if found.Carbs() != recordPfc.Carbs() {
			t.Errorf("Carbs = %v, want %v", found.Carbs(), recordPfc.Carbs())
		}
	})

	t.Run("正常系_存在しないRecordIDでnilが返る", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordPfcRepository(db)
		ctx := context.Background()

		user := testUser(t)
		eatenAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		record := testRecord(t, user.ID(), eatenAt)

		// 空のrowsを返す
		rows := sqlmock.NewRows(recordPfcColumns())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_pfcs` WHERE record_id = ?")).
			WithArgs(record.ID().String(), 1).
			WillReturnRows(rows)

		// FindByRecordID実行
		found, err := repo.FindByRecordID(ctx, record.ID())
		if err != nil {
			t.Fatalf("FindByRecordID() error = %v", err)
		}
		if found != nil {
			t.Error("found recordPfc should be nil for non-existent recordID")
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordPfcRepository(db)
		ctx := context.Background()

		user := testUser(t)
		eatenAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		record := testRecord(t, user.ID(), eatenAt)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_pfcs` WHERE record_id = ?")).
			WithArgs(record.ID().String(), 1).
			WillReturnError(errors.New("db error"))

		// FindByRecordID実行
		found, err := repo.FindByRecordID(ctx, record.ID())
		if err == nil {
			t.Error("FindByRecordID() should fail with db error")
		}
		if found != nil {
			t.Error("found recordPfc should be nil on error")
		}
	})
}

// ============================================================================
// FindByRecordIDs テスト
// ============================================================================

func TestGormRecordPfcRepository_FindByRecordIDs(t *testing.T) {
	t.Run("正常系_複数RecordIDで一括取得できる", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordPfcRepository(db)
		ctx := context.Background()

		user := testUser(t)
		eatenAt1 := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		eatenAt2 := time.Date(2024, 1, 2, 18, 0, 0, 0, time.UTC)
		record1 := testRecord(t, user.ID(), eatenAt1)
		record2 := testRecord(t, user.ID(), eatenAt2)
		recordPfc1 := testRecordPfc(t, record1.ID(), 30.0, 20.0, 50.0)
		recordPfc2 := testRecordPfc(t, record2.ID(), 40.0, 25.0, 60.0)

		recordIDs := []string{record1.ID().String(), record2.ID().String()}

		// GORMのFind()は WHERE record_id IN (?, ?)
		rows := sqlmock.NewRows(recordPfcColumns()).
			AddRow(
				recordPfc1.ID().String(),
				recordPfc1.RecordID().String(),
				recordPfc1.Protein(),
				recordPfc1.Fat(),
				recordPfc1.Carbs(),
			).
			AddRow(
				recordPfc2.ID().String(),
				recordPfc2.RecordID().String(),
				recordPfc2.Protein(),
				recordPfc2.Fat(),
				recordPfc2.Carbs(),
			)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_pfcs` WHERE record_id IN (?,?)")).
			WithArgs(recordIDs[0], recordIDs[1]).
			WillReturnRows(rows)

		// FindByRecordIDs実行
		found, err := repo.FindByRecordIDs(ctx, []vo.RecordID{record1.ID(), record2.ID()})
		if err != nil {
			t.Fatalf("FindByRecordIDs() error = %v", err)
		}
		if len(found) != 2 {
			t.Fatalf("expected 2 recordPfcs, got %d", len(found))
		}
		if !found[0].ID().Equals(recordPfc1.ID()) {
			t.Errorf("ID[0] = %v, want %v", found[0].ID(), recordPfc1.ID())
		}
		if !found[1].ID().Equals(recordPfc2.ID()) {
			t.Errorf("ID[1] = %v, want %v", found[1].ID(), recordPfc2.ID())
		}
	})

	t.Run("正常系_空スライスで空配列が返る", func(t *testing.T) {
		db, _ := setupMockDB(t)
		repo := gormPkg.NewGormRecordPfcRepository(db)
		ctx := context.Background()

		// 空スライスを渡す（record_pfc_repository.goで早期リターン）
		found, err := repo.FindByRecordIDs(ctx, []vo.RecordID{})
		if err != nil {
			t.Fatalf("FindByRecordIDs() error = %v", err)
		}
		if len(found) != 0 {
			t.Errorf("expected empty array, got %d recordPfcs", len(found))
		}
		// DBクエリは発行されないので、モック設定なし
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordPfcRepository(db)
		ctx := context.Background()

		user := testUser(t)
		eatenAt := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		record := testRecord(t, user.ID(), eatenAt)

		recordIDs := []string{record.ID().String()}

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `record_pfcs` WHERE record_id IN (?)")).
			WithArgs(recordIDs[0]).
			WillReturnError(errors.New("db error"))

		// FindByRecordIDs実行
		found, err := repo.FindByRecordIDs(ctx, []vo.RecordID{record.ID()})
		if err == nil {
			t.Error("FindByRecordIDs() should fail with db error")
		}
		if found != nil {
			t.Error("found recordPfcs should be nil on error")
		}
	})
}

// ============================================================================
// GetDailyPfc テスト
// ============================================================================

func TestGormRecordPfcRepository_GetDailyPfc(t *testing.T) {
	t.Run("正常系_日別PFC合計が取得できる", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordPfcRepository(db)
		ctx := context.Background()

		user := testUser(t)
		startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

		// 集計結果: total_protein, total_fat, total_carbs
		rows := sqlmock.NewRows([]string{"total_protein", "total_fat", "total_carbs"}).
			AddRow(100.0, 50.0, 200.0)

		// JOIN + SUM + COALESCE のクエリ
		mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(record_pfcs.protein), 0) as total_protein, COALESCE(SUM(record_pfcs.fat), 0) as total_fat, COALESCE(SUM(record_pfcs.carbs), 0) as total_carbs FROM `records` INNER JOIN record_pfcs ON records.id = record_pfcs.record_id WHERE")).
			WithArgs(user.ID().String(), startTime, endTime).
			WillReturnRows(rows)

		// GetDailyPfc実行
		result, err := repo.GetDailyPfc(ctx, user.ID(), startTime, endTime)
		if err != nil {
			t.Fatalf("GetDailyPfc() error = %v", err)
		}
		if result.Pfc.Protein() != 100.0 {
			t.Errorf("Protein = %v, want 100.0", result.Pfc.Protein())
		}
		if result.Pfc.Fat() != 50.0 {
			t.Errorf("Fat = %v, want 50.0", result.Pfc.Fat())
		}
		if result.Pfc.Carbs() != 200.0 {
			t.Errorf("Carbs = %v, want 200.0", result.Pfc.Carbs())
		}
	})

	t.Run("正常系_該当なしでゼロ値が返る", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordPfcRepository(db)
		ctx := context.Background()

		user := testUser(t)
		startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

		// ゼロ値の集計結果（COALESCE効果）
		rows := sqlmock.NewRows([]string{"total_protein", "total_fat", "total_carbs"}).
			AddRow(0.0, 0.0, 0.0)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(record_pfcs.protein), 0) as total_protein, COALESCE(SUM(record_pfcs.fat), 0) as total_fat, COALESCE(SUM(record_pfcs.carbs), 0) as total_carbs FROM `records` INNER JOIN record_pfcs ON records.id = record_pfcs.record_id WHERE")).
			WithArgs(user.ID().String(), startTime, endTime).
			WillReturnRows(rows)

		result, err := repo.GetDailyPfc(ctx, user.ID(), startTime, endTime)
		if err != nil {
			t.Fatalf("GetDailyPfc() error = %v", err)
		}
		if result.Pfc.Protein() != 0.0 {
			t.Errorf("Protein = %v, want 0.0", result.Pfc.Protein())
		}
		if result.Pfc.Fat() != 0.0 {
			t.Errorf("Fat = %v, want 0.0", result.Pfc.Fat())
		}
		if result.Pfc.Carbs() != 0.0 {
			t.Errorf("Carbs = %v, want 0.0", result.Pfc.Carbs())
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormRecordPfcRepository(db)
		ctx := context.Background()

		user := testUser(t)
		startTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endTime := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(SUM(record_pfcs.protein), 0) as total_protein, COALESCE(SUM(record_pfcs.fat), 0) as total_fat, COALESCE(SUM(record_pfcs.carbs), 0) as total_carbs FROM `records` INNER JOIN record_pfcs ON records.id = record_pfcs.record_id WHERE")).
			WithArgs(user.ID().String(), startTime, endTime).
			WillReturnError(errors.New("db error"))

		result, err := repo.GetDailyPfc(ctx, user.ID(), startTime, endTime)
		if err == nil {
			t.Error("GetDailyPfc() should fail with db error")
		}
		// DailyPfcはゼロ値で返る
		if result.Pfc.Protein() != 0.0 || result.Pfc.Fat() != 0.0 || result.Pfc.Carbs() != 0.0 {
			t.Error("result should be zero value on error")
		}
	})
}
