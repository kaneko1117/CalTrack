package gorm_test

import (
	"context"
	"testing"
)

// ============================================================================
// Save テスト
// ============================================================================

func TestGormUserRepository_Save(t *testing.T) {
	t.Run("正常系_ユーザーが保存される", func(t *testing.T) {
		cleanupTables(t)
		repo := newUserRepo()
		ctx := context.Background()

		user := testUser(t)

		// Save
		err := repo.Save(ctx, user)
		if err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		// FindByIDで検証
		saved, err := repo.FindByID(ctx, user.ID())
		if err != nil {
			t.Fatalf("FindByID() error = %v", err)
		}
		if saved == nil {
			t.Fatal("saved user should not be nil")
		}

		// 全フィールド検証
		if !saved.ID().Equals(user.ID()) {
			t.Errorf("ID = %v, want %v", saved.ID(), user.ID())
		}
		if saved.Email().String() != user.Email().String() {
			t.Errorf("Email = %v, want %v", saved.Email(), user.Email())
		}
		if saved.HashedPassword().String() != user.HashedPassword().String() {
			t.Errorf("HashedPassword = %v, want %v", saved.HashedPassword(), user.HashedPassword())
		}
		if saved.Nickname().String() != user.Nickname().String() {
			t.Errorf("Nickname = %v, want %v", saved.Nickname(), user.Nickname())
		}
		if saved.Weight().Kg() != user.Weight().Kg() {
			t.Errorf("Weight = %v, want %v", saved.Weight(), user.Weight())
		}
		if saved.Height().Cm() != user.Height().Cm() {
			t.Errorf("Height = %v, want %v", saved.Height(), user.Height())
		}
		if saved.Gender().String() != user.Gender().String() {
			t.Errorf("Gender = %v, want %v", saved.Gender(), user.Gender())
		}
		if saved.ActivityLevel().String() != user.ActivityLevel().String() {
			t.Errorf("ActivityLevel = %v, want %v", saved.ActivityLevel(), user.ActivityLevel())
		}
		if !saved.BirthDate().Time().Equal(user.BirthDate().Time()) {
			t.Errorf("BirthDate = %v, want %v", saved.BirthDate(), user.BirthDate())
		}
	})

	t.Run("異常系_重複メールで保存失敗", func(t *testing.T) {
		cleanupTables(t)
		repo := newUserRepo()
		ctx := context.Background()

		user1 := testUserWithEmail(t, "duplicate@example.com")
		user2 := testUserWithEmail(t, "duplicate@example.com")

		// 1回目は成功
		err := repo.Save(ctx, user1)
		if err != nil {
			t.Fatalf("Save() first user error = %v", err)
		}

		// 2回目は失敗（重複メール）
		err = repo.Save(ctx, user2)
		if err == nil {
			t.Error("Save() should fail with duplicate email")
		}
	})
}

// ============================================================================
// FindByEmail テスト
// ============================================================================

func TestGormUserRepository_FindByEmail(t *testing.T) {
	t.Run("正常系_メールでユーザーが見つかる", func(t *testing.T) {
		cleanupTables(t)
		repo := newUserRepo()
		ctx := context.Background()

		user := testUserWithEmail(t, "find@example.com")
		if err := repo.Save(ctx, user); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		// FindByEmail
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
		cleanupTables(t)
		repo := newUserRepo()
		ctx := context.Background()

		email := testUserWithEmail(t, "notfound@example.com").Email()

		found, err := repo.FindByEmail(ctx, email)
		if err != nil {
			t.Fatalf("FindByEmail() error = %v", err)
		}
		if found != nil {
			t.Error("found user should be nil for non-existent email")
		}
	})
}

// ============================================================================
// ExistsByEmail テスト
// ============================================================================

func TestGormUserRepository_ExistsByEmail(t *testing.T) {
	t.Run("正常系_存在するメールでtrueが返る", func(t *testing.T) {
		cleanupTables(t)
		repo := newUserRepo()
		ctx := context.Background()

		user := testUserWithEmail(t, "exists@example.com")
		if err := repo.Save(ctx, user); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		exists, err := repo.ExistsByEmail(ctx, user.Email())
		if err != nil {
			t.Fatalf("ExistsByEmail() error = %v", err)
		}
		if !exists {
			t.Error("ExistsByEmail() should return true for existing email")
		}
	})

	t.Run("正常系_存在しないメールでfalseが返る", func(t *testing.T) {
		cleanupTables(t)
		repo := newUserRepo()
		ctx := context.Background()

		email := testUserWithEmail(t, "notexists@example.com").Email()

		exists, err := repo.ExistsByEmail(ctx, email)
		if err != nil {
			t.Fatalf("ExistsByEmail() error = %v", err)
		}
		if exists {
			t.Error("ExistsByEmail() should return false for non-existent email")
		}
	})
}

// ============================================================================
// FindByID テスト
// ============================================================================

func TestGormUserRepository_FindByID(t *testing.T) {
	t.Run("正常系_IDでユーザーが見つかる", func(t *testing.T) {
		cleanupTables(t)
		repo := newUserRepo()
		ctx := context.Background()

		user := testUser(t)
		if err := repo.Save(ctx, user); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		// FindByID
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
		cleanupTables(t)
		repo := newUserRepo()
		ctx := context.Background()

		nonExistentID := testUser(t).ID()

		found, err := repo.FindByID(ctx, nonExistentID)
		if err != nil {
			t.Fatalf("FindByID() error = %v", err)
		}
		if found != nil {
			t.Error("found user should be nil for non-existent ID")
		}
	})
}

// ============================================================================
// Update テスト
// ============================================================================

func TestGormUserRepository_Update(t *testing.T) {
	t.Run("正常系_ユーザーが更新される", func(t *testing.T) {
		cleanupTables(t)
		repo := newUserRepo()
		ctx := context.Background()

		user := testUser(t)
		if err := repo.Save(ctx, user); err != nil {
			t.Fatalf("Save() error = %v", err)
		}

		// プロフィール更新
		errs := user.UpdateProfile("updatednick", 180.0, 75.0, "active")
		if errs != nil {
			t.Fatalf("UpdateProfile() errors = %v", errs)
		}

		// Update
		err := repo.Update(ctx, user)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		// FindByIDで更新後の値を検証
		updated, err := repo.FindByID(ctx, user.ID())
		if err != nil {
			t.Fatalf("FindByID() error = %v", err)
		}
		if updated == nil {
			t.Fatal("updated user should not be nil")
		}

		if updated.Nickname().String() != "updatednick" {
			t.Errorf("Nickname = %v, want updatednick", updated.Nickname())
		}
		if updated.Height().Cm() != 180.0 {
			t.Errorf("Height = %v, want 180.0", updated.Height())
		}
		if updated.Weight().Kg() != 75.0 {
			t.Errorf("Weight = %v, want 75.0", updated.Weight())
		}
		if updated.ActivityLevel().String() != "active" {
			t.Errorf("ActivityLevel = %v, want active", updated.ActivityLevel())
		}
	})
}
