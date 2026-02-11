package gorm_test

import (
	"context"
	"database/sql/driver"
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

func TestGormUserRepository_Save(t *testing.T) {
	t.Run("正常系_ユーザーが保存される", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		user := testUser(t)

		// GORMのCreate()はINSERTを実行
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
			WithArgs(
				user.ID().String(),
				user.Email().String(),
				user.HashedPassword().String(),
				user.Nickname().String(),
				user.Weight().Kg(),
				user.Height().Cm(),
				user.BirthDate().Time(),
				user.Gender().String(),
				user.ActivityLevel().String(),
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
			).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		// Save実行
		err := repo.Save(ctx, user)
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}
	})

	t.Run("異常系_DBエラーで保存失敗", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		user := testUser(t)

		// INSERT失敗をシミュレート
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO `users`")).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		// Save実行
		err := repo.Save(ctx, user)
		if err == nil {
			t.Error("Save() should fail with db error")
		}
	})
}

// ============================================================================
// FindByEmail テスト
// ============================================================================

func TestGormUserRepository_FindByEmail(t *testing.T) {
	t.Run("正常系_メールでユーザーが見つかる", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		user := testUserWithEmail(t, "find@example.com")

		// GORMのFirst()は ORDER BY id LIMIT 1 を付加
		rows := sqlmock.NewRows(userColumns()).
			AddRow(
				user.ID().String(),
				user.Email().String(),
				user.HashedPassword().String(),
				user.Nickname().String(),
				user.Weight().Kg(),
				user.Height().Cm(),
				user.BirthDate().Time(),
				user.Gender().String(),
				user.ActivityLevel().String(),
				user.CreatedAt(),
				user.UpdatedAt(),
			)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ?")).
			WithArgs(user.Email().String(), 1). // First()は最後に1を追加
			WillReturnRows(rows)

		// FindByEmail実行
		found, err := repo.FindByEmail(ctx, user.Email())
		if err != nil {
			t.Fatalf("FindByEmail() error = %v", err)
		}
		if found == nil {
			t.Fatal("found user should not be nil")
		}
		if !found.ID().Equals(user.ID()) {
			t.Errorf("ID = %v, want %v", found.ID(), user.ID())
		}
		if found.Email().String() != user.Email().String() {
			t.Errorf("Email = %v, want %v", found.Email(), user.Email())
		}
	})

	t.Run("正常系_存在しないメールでnilが返る", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		email := testUserWithEmail(t, "notfound@example.com").Email()

		// 空のrowsを返す
		rows := sqlmock.NewRows(userColumns())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ?")).
			WithArgs(email.String(), 1).
			WillReturnRows(rows)

		// FindByEmail実行
		found, err := repo.FindByEmail(ctx, email)
		if err != nil {
			t.Fatalf("FindByEmail() error = %v", err)
		}
		if found != nil {
			t.Error("found user should be nil for non-existent email")
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		email := testUserWithEmail(t, "error@example.com").Email()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE email = ?")).
			WithArgs(email.String(), 1).
			WillReturnError(errors.New("db error"))

		// FindByEmail実行
		found, err := repo.FindByEmail(ctx, email)
		if err == nil {
			t.Error("FindByEmail() should fail with db error")
		}
		if found != nil {
			t.Error("found user should be nil on error")
		}
	})
}

// ============================================================================
// ExistsByEmail テスト
// ============================================================================

func TestGormUserRepository_ExistsByEmail(t *testing.T) {
	t.Run("正常系_存在するメールでtrueが返る", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		email := testUserWithEmail(t, "exists@example.com").Email()

		// COUNT(*)クエリ
		rows := sqlmock.NewRows([]string{"count"}).AddRow(1)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE email = ?")).
			WithArgs(email.String()).
			WillReturnRows(rows)

		// ExistsByEmail実行
		exists, err := repo.ExistsByEmail(ctx, email)
		if err != nil {
			t.Fatalf("ExistsByEmail() error = %v", err)
		}
		if !exists {
			t.Error("ExistsByEmail() should return true for existing email")
		}
	})

	t.Run("正常系_存在しないメールでfalseが返る", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		email := testUserWithEmail(t, "notexists@example.com").Email()

		// COUNT(*)クエリ（結果0）
		rows := sqlmock.NewRows([]string{"count"}).AddRow(0)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE email = ?")).
			WithArgs(email.String()).
			WillReturnRows(rows)

		// ExistsByEmail実行
		exists, err := repo.ExistsByEmail(ctx, email)
		if err != nil {
			t.Fatalf("ExistsByEmail() error = %v", err)
		}
		if exists {
			t.Error("ExistsByEmail() should return false for non-existent email")
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		email := testUserWithEmail(t, "error@example.com").Email()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT count(*) FROM `users` WHERE email = ?")).
			WithArgs(email.String()).
			WillReturnError(errors.New("db error"))

		// ExistsByEmail実行
		exists, err := repo.ExistsByEmail(ctx, email)
		if err == nil {
			t.Error("ExistsByEmail() should fail with db error")
		}
		if exists {
			t.Error("ExistsByEmail() should return false on error")
		}
	})
}

// ============================================================================
// FindByID テスト
// ============================================================================

func TestGormUserRepository_FindByID(t *testing.T) {
	t.Run("正常系_IDでユーザーが見つかる", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		user := testUser(t)

		// GORMのFirst()は ORDER BY id LIMIT 1 を付加
		rows := sqlmock.NewRows(userColumns()).
			AddRow(
				user.ID().String(),
				user.Email().String(),
				user.HashedPassword().String(),
				user.Nickname().String(),
				user.Weight().Kg(),
				user.Height().Cm(),
				user.BirthDate().Time(),
				user.Gender().String(),
				user.ActivityLevel().String(),
				user.CreatedAt(),
				user.UpdatedAt(),
			)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ?")).
			WithArgs(user.ID().String(), 1). // First()は最後に1を追加
			WillReturnRows(rows)

		// FindByID実行
		found, err := repo.FindByID(ctx, user.ID())
		if err != nil {
			t.Fatalf("FindByID() error = %v", err)
		}
		if found == nil {
			t.Fatal("found user should not be nil")
		}
		if !found.ID().Equals(user.ID()) {
			t.Errorf("ID = %v, want %v", found.ID(), user.ID())
		}
	})

	t.Run("正常系_存在しないIDでnilが返る", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		nonExistentID := testUser(t).ID()

		// 空のrowsを返す
		rows := sqlmock.NewRows(userColumns())

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ?")).
			WithArgs(nonExistentID.String(), 1).
			WillReturnRows(rows)

		// FindByID実行
		found, err := repo.FindByID(ctx, nonExistentID)
		if err != nil {
			t.Fatalf("FindByID() error = %v", err)
		}
		if found != nil {
			t.Error("found user should be nil for non-existent ID")
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		id := testUser(t).ID()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE id = ?")).
			WithArgs(id.String(), 1).
			WillReturnError(errors.New("db error"))

		// FindByID実行
		found, err := repo.FindByID(ctx, id)
		if err == nil {
			t.Error("FindByID() should fail with db error")
		}
		if found != nil {
			t.Error("found user should be nil on error")
		}
	})
}

// ============================================================================
// Update テスト
// ============================================================================

func TestGormUserRepository_Update(t *testing.T) {
	t.Run("正常系_ユーザーが更新される", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		user := testUser(t)

		// GORMのSave()は UPDATE文を実行
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET")).
			WithArgs(
				user.Email().String(),
				user.HashedPassword().String(),
				user.Nickname().String(),
				user.Weight().Kg(),
				user.Height().Cm(),
				user.BirthDate().Time(),
				user.Gender().String(),
				user.ActivityLevel().String(),
				sqlmock.AnyArg(), // created_at
				sqlmock.AnyArg(), // updated_at
				user.ID().String(), // WHERE id = ?
			).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		// Update実行
		err := repo.Update(ctx, user)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}
	})

	t.Run("異常系_DBエラー", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := gormPkg.NewGormUserRepository(db)
		ctx := context.Background()

		user := testUser(t)

		// UPDATE失敗をシミュレート
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("UPDATE `users` SET")).
			WillReturnError(errors.New("db error"))
		mock.ExpectRollback()

		// Update実行
		err := repo.Update(ctx, user)
		if err == nil {
			t.Error("Update() should fail with db error")
		}
	})
}

// ============================================================================
// ヘルパー: driver.Valuer実装（時刻をtime.Timeとして扱う）
// ============================================================================

type anyTime struct{}

func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
