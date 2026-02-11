package gorm_test

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	gormPkg "caltrack/infrastructure/persistence/gorm"
)

// ============================================================================
// Save テスト
// ============================================================================

func TestGormSessionRepository_Save(t *testing.T) {
	t.Run("正常系_セッションが保存される", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormSessionRepository(db)
		ctx := context.Background()

		user := testUser(t)
		session := testSession(t, user.ID())

		// GORMのCreate()はINSERTを実行
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `sessions`")).
			WithArgs(
				session.ID().String(),
				session.UserID().String(),
				session.ExpiresAt().Time(),
				sqlmock.AnyArg(), // created_at
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Save実行
		err := repo.Save(ctx, session)
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	})

	t.Run("異常系_DBエラーで保存失敗", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormSessionRepository(db)
		ctx := context.Background()

		user := testUser(t)
		session := testSession(t, user.ID())

		// INSERT失敗をシミュレート
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `sessions`")).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		// Save実行
		err := repo.Save(ctx, session)
		if err == nil {
			t.Error("Save() should fail with db error")
		}
	})
}

// ============================================================================
// FindByID テスト
// ============================================================================

func TestGormSessionRepository_FindByID(t *testing.T) {
	t.Run("正常系_IDでセッションが見つかる", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormSessionRepository(db)
		ctx := context.Background()

		user := testUser(t)
		session := testSession(t, user.ID())

		// GORMのFirst()は ORDER BY id LIMIT 1 を付加
		rows := sqlmock.NewRows(sessionColumns()).
			AddRow(
				session.ID().String(),
				session.UserID().String(),
				session.ExpiresAt().Time(),
				session.CreatedAt(),
			)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sessions` WHERE id = ?")).
			WithArgs(session.ID().String(), 1). // First()は最後に1を追加
			WillReturnRows(rows)

		// FindByID実行
		found, err := repo.FindByID(ctx, session.ID())
		if err != nil {
			t.Fatalf("FindByID() error = %v", err)
		}
		if found == nil {
			t.Fatal("found session should not be nil")
		}
		if !found.ID().Equals(session.ID()) {
			t.Errorf("ID = %v, want %v", found.ID(), session.ID())
		}
		if !found.UserID().Equals(session.UserID()) {
			t.Errorf("UserID = %v, want %v", found.UserID(), session.UserID())
		}
	})

	t.Run("正常系_存在しないIDでnilが返る", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormSessionRepository(db)
		ctx := context.Background()

		user := testUser(t)
		nonExistentSession := testSession(t, user.ID())

		// 空のrowsを返す
		rows := sqlmock.NewRows(sessionColumns())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sessions` WHERE id = ?")).
			WithArgs(nonExistentSession.ID().String(), 1).
			WillReturnRows(rows)

		// FindByID実行
		found, err := repo.FindByID(ctx, nonExistentSession.ID())
		if err != nil {
			t.Fatalf("FindByID() error = %v", err)
		}
		if found != nil {
			t.Error("found session should be nil for non-existent ID")
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormSessionRepository(db)
		ctx := context.Background()

		user := testUser(t)
		session := testSession(t, user.ID())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `sessions` WHERE id = ?")).
			WithArgs(session.ID().String(), 1).
			WillReturnError(errors.New("db error"))

		// FindByID実行
		found, err := repo.FindByID(ctx, session.ID())
		if err == nil {
			t.Error("FindByID() should fail with db error")
		}
		if found != nil {
			t.Error("found session should be nil on error")
		}
	})
}

// ============================================================================
// DeleteByID テスト
// ============================================================================

func TestGormSessionRepository_DeleteByID(t *testing.T) {
	t.Run("正常系_セッションが削除される", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormSessionRepository(db)
		ctx := context.Background()

		user := testUser(t)
		session := testSession(t, user.ID())

		// GORMのDelete()はハードデリート（DeletedAtフィールドなし）
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `sessions` WHERE id = ?")).
			WithArgs(session.ID().String()).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// DeleteByID実行
		err := repo.DeleteByID(ctx, session.ID())
		if err != nil {
			t.Fatalf("DeleteByID() error = %v", err)
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormSessionRepository(db)
		ctx := context.Background()

		user := testUser(t)
		session := testSession(t, user.ID())

		// DELETE失敗をシミュレート
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `sessions` WHERE id = ?")).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		// DeleteByID実行
		err := repo.DeleteByID(ctx, session.ID())
		if err == nil {
			t.Error("DeleteByID() should fail with db error")
		}
	})
}

// ============================================================================
// DeleteByUserID テスト
// ============================================================================

func TestGormSessionRepository_DeleteByUserID(t *testing.T) {
	t.Run("正常系_ユーザーの全セッションが削除される", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormSessionRepository(db)
		ctx := context.Background()

		user := testUser(t)

		// GORMのDelete()はハードデリート
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `sessions` WHERE user_id = ?")).
			WithArgs(user.ID().String()).
			WillReturnResult(sqlmock.NewResult(0, 2)) // 2件削除されたと想定
		mock.ExpectCommit()

		// DeleteByUserID実行
		err := repo.DeleteByUserID(ctx, user.ID())
		if err != nil {
			t.Fatalf("DeleteByUserID() error = %v", err)
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormSessionRepository(db)
		ctx := context.Background()

		user := testUser(t)

		// DELETE失敗をシミュレート
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM `sessions` WHERE user_id = ?")).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		// DeleteByUserID実行
		err := repo.DeleteByUserID(ctx, user.ID())
		if err == nil {
			t.Error("DeleteByUserID() should fail with db error")
		}
	})
}
