package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/mock"
	"caltrack/usecase"

	gomock "go.uber.org/mock/gomock"
)

// setupUserMocks はUser Usecase用のモックを初期化する
func setupUserMocks(t *testing.T) (*mock.MockUserRepository, *mock.MockTransactionManager, *gomock.Controller) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := mock.NewMockUserRepository(ctrl)
	txManager := mock.NewMockTransactionManager(ctrl)
	return userRepo, txManager, ctrl
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
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		user := validUser(t)
		email, _ := vo.NewEmail("test@example.com")

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			ExistsByEmail(gomock.Any(), gomock.Eq(email)).
			Return(false, nil)
		userRepo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(nil)

		uc := usecase.NewUserUsecase(userRepo, txManager)
		registeredUser, err := uc.Register(context.Background(), user)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if registeredUser.ID().String() == "" {
			t.Error("UserID should not be empty")
		}
	})

	t.Run("異常系_メールアドレスが既に存在する", func(t *testing.T) {
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		user := validUser(t)
		email, _ := vo.NewEmail("test@example.com")

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			ExistsByEmail(gomock.Any(), gomock.Eq(email)).
			Return(true, nil)

		uc := usecase.NewUserUsecase(userRepo, txManager)
		_, err := uc.Register(context.Background(), user)

		if err != domainErrors.ErrEmailAlreadyExists {
			t.Errorf("got %v, want ErrEmailAlreadyExists", err)
		}
	})

	t.Run("異常系_リポジトリエラー", func(t *testing.T) {
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		user := validUser(t)
		email, _ := vo.NewEmail("test@example.com")
		repoErr := errors.New("db error")

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			ExistsByEmail(gomock.Any(), gomock.Eq(email)).
			Return(false, repoErr)

		uc := usecase.NewUserUsecase(userRepo, txManager)
		_, err := uc.Register(context.Background(), user)

		if err != repoErr {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_保存エラー", func(t *testing.T) {
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		user := validUser(t)
		email, _ := vo.NewEmail("test@example.com")
		saveErr := errors.New("save error")

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			ExistsByEmail(gomock.Any(), gomock.Eq(email)).
			Return(false, nil)
		userRepo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(saveErr)

		uc := usecase.NewUserUsecase(userRepo, txManager)
		_, err := uc.Register(context.Background(), user)

		if !errors.Is(err, saveErr) {
			t.Errorf("got %v, want saveErr", err)
		}
	})
}

// TestUserUsecase_UpdateProfile はプロフィール更新機能のテスト
func TestUserUsecase_UpdateProfile(t *testing.T) {
	t.Run("正常系_プロフィール更新成功", func(t *testing.T) {
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		user := reconstructedUser(t)

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(user.ID())).
			Return(user, nil)
		userRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(nil)

		uc := usecase.NewUserUsecase(userRepo, txManager)
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
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(nil, nil)

		uc := usecase.NewUserUsecase(userRepo, txManager)
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
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		user := reconstructedUser(t)

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(user.ID())).
			Return(user, nil)

		uc := usecase.NewUserUsecase(userRepo, txManager)
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
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		repoErr := errors.New("db error")

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(nil, repoErr)

		uc := usecase.NewUserUsecase(userRepo, txManager)
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
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		user := reconstructedUser(t)
		updateErr := errors.New("update error")

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(user.ID())).
			Return(user, nil)
		userRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(updateErr)

		uc := usecase.NewUserUsecase(userRepo, txManager)
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
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		user := reconstructedUser(t)

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(user.ID())).
			Return(user, nil)
		userRepo.EXPECT().
			Update(gomock.Any(), gomock.Any()).
			Return(nil)

		uc := usecase.NewUserUsecase(userRepo, txManager)
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
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		user := reconstructedUser(t)

		// GetProfileはトランザクションを使用しないのでEXPECTは不要
		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(user.ID())).
			Return(user, nil)

		uc := usecase.NewUserUsecase(userRepo, txManager)
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
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(nil, nil)

		uc := usecase.NewUserUsecase(userRepo, txManager)
		_, err := uc.GetProfile(context.Background(), userID)

		if err != domainErrors.ErrUserNotFound {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_リポジトリエラー", func(t *testing.T) {
		userRepo, txManager, ctrl := setupUserMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		repoErr := errors.New("db error")

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(nil, repoErr)

		uc := usecase.NewUserUsecase(userRepo, txManager)
		_, err := uc.GetProfile(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})
}
