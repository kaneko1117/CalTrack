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
// Execute テスト
// ============================================================================

func TestGormTransactionManager_Execute(t *testing.T) {
	t.Run("正常系_トランザクション内で操作がコミットされる", func(t *testing.T) {
		db, mock := setupMockDB(t)
		tm := gormPkg.NewGormTransactionManager(db)
		ctx := context.Background()

		// GORMのTransaction()は内部でBegin→処理→Commitを実行
		mock.ExpectBegin()

		// トランザクション内での操作をシミュレート（例: INSERT）
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `test_table`")).
			WithArgs("test_value").
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectCommit()

		// Execute実行
		err := tm.Execute(ctx, func(txCtx context.Context) error {
			// トランザクション内の処理
			txDB := gormPkg.GetTx(txCtx, db)
			return txDB.Exec("INSERT INTO `test_table` (col) VALUES (?)", "test_value").Error
		})

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}
	})

	t.Run("異常系_エラー時にロールバックされる", func(t *testing.T) {
		db, mock := setupMockDB(t)
		tm := gormPkg.NewGormTransactionManager(db)
		ctx := context.Background()

		// GORMのTransaction()はエラー時にRollbackを実行
		mock.ExpectBegin()
		mock.ExpectRollback()

		testErr := errors.New("transaction error")

		// Execute実行
		err := tm.Execute(ctx, func(txCtx context.Context) error {
			// エラーを返すとロールバックされる
			return testErr
		})

		if err == nil {
			t.Error("Execute() should fail with transaction error")
		}
		if err != testErr {
			t.Errorf("Execute() error = %v, want %v", err, testErr)
		}
	})

	t.Run("異常系_Begin失敗", func(t *testing.T) {
		db, mock := setupMockDB(t)
		tm := gormPkg.NewGormTransactionManager(db)
		ctx := context.Background()

		// Begin失敗をシミュレート
		beginErr := errors.New("begin failed")
		mock.ExpectBegin().WillReturnError(beginErr)

		// Execute実行
		err := tm.Execute(ctx, func(txCtx context.Context) error {
			// この関数は実行されない
			return nil
		})

		if err == nil {
			t.Error("Execute() should fail when Begin fails")
		}
	})
}

// ============================================================================
// GetTx テスト
// ============================================================================

func TestGetTx(t *testing.T) {
	t.Run("コンテキストにtxがある場合_txが返る", func(t *testing.T) {
		db, mock := setupMockDB(t)
		tm := gormPkg.NewGormTransactionManager(db)
		ctx := context.Background()

		var capturedTxDB interface{}

		// トランザクション内でGetTxを呼び出し、txを取得
		mock.ExpectBegin()
		mock.ExpectCommit()

		err := tm.Execute(ctx, func(txCtx context.Context) error {
			txDB := gormPkg.GetTx(txCtx, db)
			capturedTxDB = txDB
			return nil
		})

		if err != nil {
			t.Fatalf("Execute() error = %v", err)
		}

		// capturedTxDBがdbと異なることを確認（txが返されている）
		if capturedTxDB == nil {
			t.Error("GetTx() should return tx from context")
		}
		if capturedTxDB == db {
			t.Error("GetTx() should return tx, not db")
		}
	})

	t.Run("コンテキストにtxがない場合_dbが返る", func(t *testing.T) {
		db, _ := setupMockDB(t)
		ctx := context.Background()

		// トランザクション外でGetTxを呼び出し
		result := gormPkg.GetTx(ctx, db)

		// dbが返されることを確認
		if result != db {
			t.Error("GetTx() should return db when no tx in context")
		}
	})
}
