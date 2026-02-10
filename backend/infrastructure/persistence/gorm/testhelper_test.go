package gorm_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"

	"caltrack/domain/entity"
	"caltrack/domain/vo"
	gormPkg "caltrack/infrastructure/persistence/gorm"
	"caltrack/infrastructure/persistence/gorm/model"
)

var testDB *gorm.DB

// TestMain は全テストの前後処理を行う
func TestMain(m *testing.M) {
	// SQLite in-memoryでDB接続
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open sqlite in-memory: %v", err)
	}
	testDB = db

	// AutoMigrateで全モデルのテーブル作成
	if err := testDB.AutoMigrate(
		&model.User{},
		&model.Record{},
		&model.RecordItem{},
		&model.RecordPfc{},
		&model.Session{},
		&model.AdviceCache{},
	); err != nil {
		log.Fatalf("failed to auto migrate: %v", err)
	}

	// テスト実行
	code := m.Run()
	os.Exit(code)
}

// cleanupTables は全テーブルを削除する（SQLite用）
func cleanupTables(t *testing.T) {
	t.Helper()

	// 外部キー制約を無視して全レコードを削除
	tables := []string{
		"record_items",
		"record_pfcs",
		"records",
		"sessions",
		"advice_caches",
		"users",
	}
	for _, table := range tables {
		if err := testDB.Exec("DELETE FROM " + table).Error; err != nil {
			t.Fatalf("failed to delete from %s: %v", table, err)
		}
	}
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
// Repository Helpers
// ============================================================================

// newUserRepo は新しいGormUserRepositoryを生成する
func newUserRepo() *gormPkg.GormUserRepository {
	return gormPkg.NewGormUserRepository(testDB)
}

// newRecordRepo は新しいGormRecordRepositoryを生成する
func newRecordRepo() *gormPkg.GormRecordRepository {
	return gormPkg.NewGormRecordRepository(testDB)
}

// newRecordPfcRepo は新しいGormRecordPfcRepositoryを生成する
func newRecordPfcRepo() *gormPkg.GormRecordPfcRepository {
	return gormPkg.NewGormRecordPfcRepository(testDB)
}

// newSessionRepo は新しいGormSessionRepositoryを生成する
func newSessionRepo() *gormPkg.GormSessionRepository {
	return gormPkg.NewGormSessionRepository(testDB)
}

// newAdviceCacheRepo は新しいGormAdviceCacheRepositoryを生成する
func newAdviceCacheRepo() *gormPkg.GormAdviceCacheRepository {
	return gormPkg.NewGormAdviceCacheRepository(testDB)
}

// newTxManager は新しいGormTransactionManagerを生成する
func newTxManager() *gormPkg.GormTransactionManager {
	return gormPkg.NewGormTransactionManager(testDB)
}

// ============================================================================
// Data Setup Helpers
// ============================================================================

// saveTestUser はテスト用ユーザーをDBに保存する
func saveTestUser(t *testing.T, user *entity.User) {
	t.Helper()
	repo := newUserRepo()
	if err := repo.Save(context.Background(), user); err != nil {
		t.Fatalf("failed to save test user: %v", err)
	}
}

// saveTestRecord はテスト用RecordをDBに保存する
func saveTestRecord(t *testing.T, record *entity.Record) {
	t.Helper()
	repo := newRecordRepo()
	if err := repo.Save(context.Background(), record); err != nil {
		t.Fatalf("failed to save test record: %v", err)
	}
}

// saveTestRecordPfc はテスト用RecordPfcをDBに保存する
func saveTestRecordPfc(t *testing.T, recordPfc *entity.RecordPfc) {
	t.Helper()
	repo := newRecordPfcRepo()
	if err := repo.Save(context.Background(), recordPfc); err != nil {
		t.Fatalf("failed to save test record pfc: %v", err)
	}
}

// saveTestSession はテスト用SessionをDBに保存する
func saveTestSession(t *testing.T, session *entity.Session) {
	t.Helper()
	repo := newSessionRepo()
	if err := repo.Save(context.Background(), session); err != nil {
		t.Fatalf("failed to save test session: %v", err)
	}
}

// saveTestAdviceCache はテスト用AdviceCacheをDBに保存する
func saveTestAdviceCache(t *testing.T, cache *entity.AdviceCache) {
	t.Helper()
	repo := newAdviceCacheRepo()
	if err := repo.Save(context.Background(), cache); err != nil {
		t.Fatalf("failed to save test advice cache: %v", err)
	}
}
