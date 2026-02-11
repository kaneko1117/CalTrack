package gorm_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
)

// setupMockDB はgo-sqlmockでモックDBを作成し、GORMのDBインスタンスを返す
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}

	// mysqlダイアレクタでGORMを初期化
	dialector := mysql.New(mysql.Config{
		Conn:                      sqlDB,
		SkipInitializeWithVersion: true,
	})

	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	// クリーンアップ: 期待値が全て消費されたか確認
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unfulfilled expectations: %v", err)
		}
		sqlDB.Close()
	})

	return db, mock
}

// ============================================================================
// Entity Factory Helpers
// ============================================================================

// testUser はデフォルトのテスト用ユーザーを生成する
func testUser(t *testing.T) *entity.User {
	t.Helper()
	return testUserWithEmail(t, "test@example.com")
}

// testUserWithEmail はメールアドレスを指定してテスト用ユーザーを生成する
func testUserWithEmail(t *testing.T, email string) *entity.User {
	t.Helper()
	user, errs := entity.NewUser(
		email,
		"password123",
		"testuser",
		70.0,
		175.0,
		time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		"male",
		"moderate",
	)
	if errs != nil {
		t.Fatalf("failed to create test user: %v", errs)
	}
	return user
}

// testRecord はテスト用Recordを生成する
func testRecord(t *testing.T, userID vo.UserID, eatenAt time.Time) *entity.Record {
	t.Helper()
	record, err := entity.NewRecord(userID, eatenAt)
	if err != nil {
		t.Fatalf("failed to create test record: %v", err)
	}
	return record
}

// testRecordWithItem は食品名とカロリーを持つテスト用Recordを生成する
func testRecordWithItem(t *testing.T, userID vo.UserID, eatenAt time.Time, itemName string, calories int) *entity.Record {
	t.Helper()
	record := testRecord(t, userID, eatenAt)
	if err := record.AddItem(itemName, calories); err != nil {
		t.Fatalf("failed to add item to record: %v", err)
	}
	return record
}

// testSession はテスト用Sessionを生成する
func testSession(t *testing.T, userID vo.UserID) *entity.Session {
	t.Helper()
	session, err := entity.NewSessionWithUserID(userID)
	if err != nil {
		t.Fatalf("failed to create test session: %v", err)
	}
	return session
}

// testAdviceCache はテスト用AdviceCacheを生成する
func testAdviceCache(t *testing.T, userID vo.UserID, date time.Time, advice string) *entity.AdviceCache {
	t.Helper()
	return entity.NewAdviceCache(userID, date, advice)
}

// testRecordPfc はテスト用RecordPfcを生成する
func testRecordPfc(t *testing.T, recordID vo.RecordID, protein, fat, carbs float64) *entity.RecordPfc {
	t.Helper()
	return entity.NewRecordPfc(recordID, protein, fat, carbs)
}

// ============================================================================
// カラム定義ヘルパー
// ============================================================================

// userColumns はUsersテーブルのカラム一覧を返す
func userColumns() []string {
	return []string{
		"id",
		"email",
		"hashed_password",
		"nickname",
		"weight",
		"height",
		"birth_date",
		"gender",
		"activity_level",
		"created_at",
		"updated_at",
	}
}
