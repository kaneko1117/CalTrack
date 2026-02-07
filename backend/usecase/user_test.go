package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/usecase"
)

type mockUserRepository struct {
	existsByEmail func(ctx context.Context, email vo.Email) (bool, error)
	save          func(ctx context.Context, user *entity.User) error
	findByID      func(ctx context.Context, id vo.UserID) (*entity.User, error)
	update        func(ctx context.Context, user *entity.User) error
}

func (m *mockUserRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	return m.existsByEmail(ctx, email)
}

func (m *mockUserRepository) Save(ctx context.Context, u *entity.User) error {
	return m.save(ctx, u)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	return nil, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id vo.UserID) (*entity.User, error) {
	if m.findByID != nil {
		return m.findByID(ctx, id)
	}
	return nil, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *entity.User) error {
	if m.update != nil {
		return m.update(ctx, user)
	}
	return nil
}

type mockTransactionManager struct{}

func (m *mockTransactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func validUser(t *testing.T) *entity.User {
	t.Helper()
	u, errs := entity.NewUser(
		"test@example.com",
		"password123",
		"testuser",
		70.5,
		175.0,
		time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		"male",
		"moderate",
	)
	if errs != nil {
		t.Fatalf("failed to create valid user: %v", errs)
	}
	return u
}

func reconstructedUser(t *testing.T) *entity.User {
	t.Helper()
	u, err := entity.ReconstructUser(
		"550e8400-e29b-41d4-a716-446655440000",
		"test@example.com",
		"$2a$10$hashedpassword",
		"oldnick",
		60.0,
		165.0,
		time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		"male",
		"sedentary",
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	)
	if err != nil {
		t.Fatalf("failed to reconstruct user: %v", err)
	}
	return u
}

// TestUserUsecase_Register はユーザー登録機能のテスト
func TestUserUsecase_Register(t *testing.T) {
	t.Run("正常系_登録成功", func(t *testing.T) {
		repo := &mockUserRepository{
			existsByEmail: func(ctx context.Context, email vo.Email) (bool, error) {
				return false, nil
			},
			save: func(ctx context.Context, u *entity.User) error {
				return nil
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		registeredUser, err := uc.Register(context.Background(), validUser(t))

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if registeredUser.ID().String() == "" {
			t.Error("UserID should not be empty")
		}
	})

	t.Run("異常系_メールアドレスが既に存在する", func(t *testing.T) {
		repo := &mockUserRepository{
			existsByEmail: func(ctx context.Context, email vo.Email) (bool, error) {
				return true, nil
			},
			save: func(ctx context.Context, u *entity.User) error {
				return nil
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		_, err := uc.Register(context.Background(), validUser(t))

		if err != domainErrors.ErrEmailAlreadyExists {
			t.Errorf("got %v, want ErrEmailAlreadyExists", err)
		}
	})

	t.Run("異常系_リポジトリエラー", func(t *testing.T) {
		repoErr := errors.New("db error")
		repo := &mockUserRepository{
			existsByEmail: func(ctx context.Context, email vo.Email) (bool, error) {
				return false, repoErr
			},
			save: func(ctx context.Context, u *entity.User) error {
				return nil
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		_, err := uc.Register(context.Background(), validUser(t))

		if err != repoErr {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_保存エラー", func(t *testing.T) {
		saveErr := errors.New("save error")
		repo := &mockUserRepository{
			existsByEmail: func(ctx context.Context, email vo.Email) (bool, error) {
				return false, nil
			},
			save: func(ctx context.Context, u *entity.User) error {
				return saveErr
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		_, err := uc.Register(context.Background(), validUser(t))

		if !errors.Is(err, saveErr) {
			t.Errorf("got %v, want saveErr", err)
		}
	})
}

// TestUserUsecase_UpdateProfile はプロフィール更新機能のテスト
func TestUserUsecase_UpdateProfile(t *testing.T) {
	t.Run("正常系_プロフィール更新成功", func(t *testing.T) {
		user := reconstructedUser(t)
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
			update: func(ctx context.Context, u *entity.User) error {
				return nil
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		updatedUser, err := uc.UpdateProfile(context.Background(), user.ID(), usecase.UpdateProfileInput{
			Nickname:      "newnick",
			Height:        170.0,
			Weight:        65.0,
			ActivityLevel: "moderate",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updatedUser.Nickname().String() != "newnick" {
			t.Errorf("nickname got %v, want newnick", updatedUser.Nickname().String())
		}
		if updatedUser.Height().Cm() != 170.0 {
			t.Errorf("height got %v, want 170.0", updatedUser.Height().Cm())
		}
		if updatedUser.Weight().Kg() != 65.0 {
			t.Errorf("weight got %v, want 65.0", updatedUser.Weight().Kg())
		}
		if updatedUser.ActivityLevel().String() != "moderate" {
			t.Errorf("activity level got %v, want moderate", updatedUser.ActivityLevel().String())
		}
	})

	t.Run("異常系_ユーザーが見つからない", func(t *testing.T) {
		userID := vo.NewUserID()
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, nil
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		_, err := uc.UpdateProfile(context.Background(), userID, usecase.UpdateProfileInput{
			Nickname:      "newnick",
			Height:        170.0,
			Weight:        65.0,
			ActivityLevel: "moderate",
		})

		if err != domainErrors.ErrUserNotFound {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_バリデーションエラー", func(t *testing.T) {
		user := reconstructedUser(t)
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		_, err := uc.UpdateProfile(context.Background(), user.ID(), usecase.UpdateProfileInput{
			Nickname:      "", // ニックネームが空
			Height:        170.0,
			Weight:        65.0,
			ActivityLevel: "moderate",
		})

		if err != domainErrors.ErrNicknameRequired {
			t.Errorf("got %v, want ErrNicknameRequired", err)
		}
	})

	t.Run("異常系_FindByIDリポジトリエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		repoErr := errors.New("db error")
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, repoErr
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		_, err := uc.UpdateProfile(context.Background(), userID, usecase.UpdateProfileInput{
			Nickname:      "newnick",
			Height:        170.0,
			Weight:        65.0,
			ActivityLevel: "moderate",
		})

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_Updateリポジトリエラー", func(t *testing.T) {
		user := reconstructedUser(t)
		updateErr := errors.New("update error")
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
			update: func(ctx context.Context, u *entity.User) error {
				return updateErr
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		_, err := uc.UpdateProfile(context.Background(), user.ID(), usecase.UpdateProfileInput{
			Nickname:      "newnick",
			Height:        170.0,
			Weight:        65.0,
			ActivityLevel: "moderate",
		})

		if !errors.Is(err, updateErr) {
			t.Errorf("got %v, want updateErr", err)
		}
	})

	t.Run("正常系_更新後のEntityが返却される", func(t *testing.T) {
		user := reconstructedUser(t)
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
			update: func(ctx context.Context, u *entity.User) error {
				return nil
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		updatedUser, err := uc.UpdateProfile(context.Background(), user.ID(), usecase.UpdateProfileInput{
			Nickname:      "updatednick",
			Height:        180.0,
			Weight:        75.0,
			ActivityLevel: "active",
		})

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if updatedUser == nil {
			t.Fatal("updatedUser should not be nil")
		}
		if updatedUser.ID().String() != user.ID().String() {
			t.Errorf("user ID mismatch: got %v, want %v", updatedUser.ID().String(), user.ID().String())
		}
		if updatedUser.Nickname().String() != "updatednick" {
			t.Errorf("nickname got %v, want updatednick", updatedUser.Nickname().String())
		}
	})
}

// TestUserUsecase_GetProfile はユーザー情報取得機能のテスト
func TestUserUsecase_GetProfile(t *testing.T) {
	t.Run("正常系_プロフィール取得成功", func(t *testing.T) {
		user := reconstructedUser(t)
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		result, err := uc.GetProfile(context.Background(), user.ID())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("result should not be nil")
		}
		if result.Email().String() != "test@example.com" {
			t.Errorf("email got %v, want test@example.com", result.Email().String())
		}
		if result.Nickname().String() != "oldnick" {
			t.Errorf("nickname got %v, want oldnick", result.Nickname().String())
		}
		if result.Weight().Kg() != 60.0 {
			t.Errorf("weight got %v, want 60.0", result.Weight().Kg())
		}
		if result.Height().Cm() != 165.0 {
			t.Errorf("height got %v, want 165.0", result.Height().Cm())
		}
		if result.Gender().String() != "male" {
			t.Errorf("gender got %v, want male", result.Gender().String())
		}
		if result.ActivityLevel().String() != "sedentary" {
			t.Errorf("activity level got %v, want sedentary", result.ActivityLevel().String())
		}
	})

	t.Run("異常系_ユーザーが見つからない", func(t *testing.T) {
		userID := vo.NewUserID()
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, nil
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		_, err := uc.GetProfile(context.Background(), userID)

		if err != domainErrors.ErrUserNotFound {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_リポジトリエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		repoErr := errors.New("db error")
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, repoErr
			},
		}
		txManager := &mockTransactionManager{}

		uc := usecase.NewUserUsecase(repo, txManager)
		_, err := uc.GetProfile(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})
}
