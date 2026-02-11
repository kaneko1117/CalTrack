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

func TestGormAdviceCacheRepository_Save(t *testing.T) {
	t.Run("正常系_キャッシュが保存される", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormAdviceCacheRepository(db)
		ctx := context.Background()

		user := testUser(t)
		date := time.Date(2025, 2, 11, 15, 30, 0, 0, time.UTC)
		cache := testAdviceCache(t, user.ID(), date, "今日のアドバイスです")

		// GORMのCreate()はINSERTを実行
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `advice_caches`")).
			WithArgs(
				cache.ID().String(),
				cache.UserID().String(),
				cache.CacheDate(), // 日付は正規化済み（時刻部分は0）
				cache.Advice(),
				sqlmock.AnyArg(), // created_at
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Save実行
		err := repo.Save(ctx, cache)
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormAdviceCacheRepository(db)
		ctx := context.Background()

		user := testUser(t)
		date := time.Date(2025, 2, 11, 15, 30, 0, 0, time.UTC)
		cache := testAdviceCache(t, user.ID(), date, "今日のアドバイスです")

		// INSERT失敗をシミュレート
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `advice_caches`")).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		// Save実行
		err := repo.Save(ctx, cache)
		if err == nil {
			t.Error("Save() should fail with db error")
		}
	})
}

// ============================================================================
// FindByUserIDAndDate テスト
// ============================================================================

func TestGormAdviceCacheRepository_FindByUserIDAndDate(t *testing.T) {
	t.Run("正常系_ユーザーIDと日付でキャッシュが見つかる", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormAdviceCacheRepository(db)
		ctx := context.Background()

		user := testUser(t)
		date := time.Date(2025, 2, 11, 15, 30, 0, 0, time.UTC)
		cache := testAdviceCache(t, user.ID(), date, "今日のアドバイスです")

		// 正規化された日付（時刻部分が0になる）
		normalizedDate := time.Date(2025, 2, 11, 0, 0, 0, 0, time.UTC)

		// GORMのFirst()は ORDER BY id LIMIT 1 を付加
		rows := sqlmock.NewRows(adviceCacheColumns()).
			AddRow(
				cache.ID().String(),
				cache.UserID().String(),
				normalizedDate,
				cache.Advice(),
				cache.CreatedAt(),
			)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `advice_caches` WHERE user_id = ? AND cache_date = ?")).
			WithArgs(user.ID().String(), normalizedDate, 1). // First()は最後に1を追加
			WillReturnRows(rows)

		// FindByUserIDAndDate実行
		found, err := repo.FindByUserIDAndDate(ctx, user.ID(), date)
		if err != nil {
			t.Fatalf("FindByUserIDAndDate() error = %v", err)
		}
		if found == nil {
			t.Fatal("found cache should not be nil")
		}
		if !found.ID().Equals(cache.ID()) {
			t.Errorf("ID = %v, want %v", found.ID(), cache.ID())
		}
		if found.Advice() != cache.Advice() {
			t.Errorf("Advice = %v, want %v", found.Advice(), cache.Advice())
		}
	})

	t.Run("正常系_存在しない場合にnilが返る", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormAdviceCacheRepository(db)
		ctx := context.Background()

		user := testUser(t)
		date := time.Date(2025, 2, 11, 15, 30, 0, 0, time.UTC)
		normalizedDate := time.Date(2025, 2, 11, 0, 0, 0, 0, time.UTC)

		// 空のrowsを返す
		rows := sqlmock.NewRows(adviceCacheColumns())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `advice_caches` WHERE user_id = ? AND cache_date = ?")).
			WithArgs(user.ID().String(), normalizedDate, 1).
			WillReturnRows(rows)

		// FindByUserIDAndDate実行
		found, err := repo.FindByUserIDAndDate(ctx, user.ID(), date)
		if err != nil {
			t.Fatalf("FindByUserIDAndDate() error = %v", err)
		}
		if found != nil {
			t.Error("found cache should be nil for non-existent record")
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormAdviceCacheRepository(db)
		ctx := context.Background()

		user := testUser(t)
		date := time.Date(2025, 2, 11, 15, 30, 0, 0, time.UTC)
		normalizedDate := time.Date(2025, 2, 11, 0, 0, 0, 0, time.UTC)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `advice_caches` WHERE user_id = ? AND cache_date = ?")).
			WithArgs(user.ID().String(), normalizedDate, 1).
			WillReturnError(errors.New("db error"))

		// FindByUserIDAndDate実行
		found, err := repo.FindByUserIDAndDate(ctx, user.ID(), date)
		if err == nil {
			t.Error("FindByUserIDAndDate() should fail with db error")
		}
		if found != nil {
			t.Error("found cache should be nil on error")
		}
	})
}

// ============================================================================
// DeleteByUserIDAndDate テスト
// ============================================================================

func TestGormAdviceCacheRepository_DeleteByUserIDAndDate(t *testing.T) {
	t.Run("正常系_キャッシュが削除される", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormAdviceCacheRepository(db)
		ctx := context.Background()

		user := testUser(t)
		date := time.Date(2025, 2, 11, 15, 30, 0, 0, time.UTC)
		normalizedDate := time.Date(2025, 2, 11, 0, 0, 0, 0, time.UTC)

		// GORMのDelete()はDELETEを実行
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `advice_caches` WHERE user_id = ? AND cache_date = ?")).
			WithArgs(user.ID().String(), normalizedDate).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// DeleteByUserIDAndDate実行
		err := repo.DeleteByUserIDAndDate(ctx, user.ID(), date)
		if err != nil {
			t.Fatalf("DeleteByUserIDAndDate() error = %v", err)
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormAdviceCacheRepository(db)
		ctx := context.Background()

		user := testUser(t)
		date := time.Date(2025, 2, 11, 15, 30, 0, 0, time.UTC)
		normalizedDate := time.Date(2025, 2, 11, 0, 0, 0, 0, time.UTC)

		// DELETE失敗をシミュレート
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `advice_caches` WHERE user_id = ? AND cache_date = ?")).
			WithArgs(user.ID().String(), normalizedDate).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		// DeleteByUserIDAndDate実行
		err := repo.DeleteByUserIDAndDate(ctx, user.ID(), date)
		if err == nil {
			t.Error("DeleteByUserIDAndDate() should fail with db error")
		}
	})
}
