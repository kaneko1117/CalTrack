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

// setupAuthMocks はテスト用のモックを初期化する
func setupAuthMocks(t *testing.T) (*mock.MockUserRepository, *mock.MockSessionRepository, *mock.MockTransactionManager, *gomock.Controller) {
	t.Helper()
	ctrl := gomock.NewController(t)
	userRepo := mock.NewMockUserRepository(ctrl)
	sessionRepo := mock.NewMockSessionRepository(ctrl)
	txManager := mock.NewMockTransactionManager(ctrl)
	return userRepo, sessionRepo, txManager, ctrl
}

// setupTxManagerExecute はTransactionManager.ExecuteのDoAndReturnを設定する
func setupTxManagerExecute(txManager *mock.MockTransactionManager) {
	txManager.EXPECT().
		Execute(gomock.Any(), gomock.Any()).
		DoAndReturn(func(ctx context.Context, fn func(ctx context.Context) error) error {
			return fn(ctx)
		})
}

// validUserForAuth はテスト用の有効なユーザーを生成する
func validUserForAuth(t *testing.T) *entity.User {
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

// validSessionID はテスト用の有効なセッションIDを生成する
func validSessionID(t *testing.T) vo.SessionID {
	t.Helper()
	sid, err := vo.NewSessionID()
	if err != nil {
		t.Fatalf("failed to create session id: %v", err)
	}
	return sid
}

// validSession はテスト用の有効なセッションを生成する
func validSession(t *testing.T, userID vo.UserID) *entity.Session {
	t.Helper()
	session, err := entity.NewSessionWithUserID(userID)
	if err != nil {
		t.Fatalf("failed to create session: %v", err)
	}
	return session
}

// =============================================================================
// Login テスト
// =============================================================================

// TestAuthUsecase_Login はログイン機能のテスト
func TestAuthUsecase_Login(t *testing.T) {
	t.Run("正常系_ログイン成功", func(t *testing.T) {
		userRepo, sessionRepo, txManager, ctrl := setupAuthMocks(t)
		defer ctrl.Finish()

		user := validUserForAuth(t)

		// VOを作成してgomock.Eq()で比較
		email, _ := vo.NewEmail("test@example.com")

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			FindByEmail(gomock.Any(), gomock.Eq(email)).
			Return(user, nil)
		sessionRepo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(nil)

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		testEmail, _ := vo.NewEmail("test@example.com")
		testPassword, _ := vo.NewPassword("password123")
		output, err := uc.Login(context.Background(), testEmail, testPassword)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if output == nil {
			t.Fatal("output should not be nil")
		}
		if output.Session == nil {
			t.Error("session should not be nil")
		}
		if output.User == nil {
			t.Error("user should not be nil")
		}
		if output.Session.UserID().String() != user.ID().String() {
			t.Errorf("session user id mismatch: got %v, want %v", output.Session.UserID().String(), user.ID().String())
		}
	})

	t.Run("異常系_ユーザーが見つからない", func(t *testing.T) {
		userRepo, sessionRepo, txManager, ctrl := setupAuthMocks(t)
		defer ctrl.Finish()

		email, _ := vo.NewEmail("notfound@example.com")

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			FindByEmail(gomock.Any(), gomock.Eq(email)).
			Return(nil, nil)

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		testEmail, _ := vo.NewEmail("notfound@example.com")
		testPassword, _ := vo.NewPassword("password123")
		_, err := uc.Login(context.Background(), testEmail, testPassword)

		if !errors.Is(err, domainErrors.ErrInvalidCredentials) {
			t.Errorf("got %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("異常系_パスワードが一致しない", func(t *testing.T) {
		userRepo, sessionRepo, txManager, ctrl := setupAuthMocks(t)
		defer ctrl.Finish()

		user := validUserForAuth(t)
		email, _ := vo.NewEmail("test@example.com")

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			FindByEmail(gomock.Any(), gomock.Eq(email)).
			Return(user, nil)

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		testEmail, _ := vo.NewEmail("test@example.com")
		testPassword, _ := vo.NewPassword("wrongpassword")
		_, err := uc.Login(context.Background(), testEmail, testPassword)

		if !errors.Is(err, domainErrors.ErrInvalidCredentials) {
			t.Errorf("got %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("異常系_ユーザーリポジトリエラー", func(t *testing.T) {
		userRepo, sessionRepo, txManager, ctrl := setupAuthMocks(t)
		defer ctrl.Finish()

		repoErr := errors.New("db error")
		email, _ := vo.NewEmail("test@example.com")

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			FindByEmail(gomock.Any(), gomock.Eq(email)).
			Return(nil, repoErr)

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		testEmail, _ := vo.NewEmail("test@example.com")
		testPassword, _ := vo.NewPassword("password123")
		_, err := uc.Login(context.Background(), testEmail, testPassword)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_セッション保存エラー", func(t *testing.T) {
		userRepo, sessionRepo, txManager, ctrl := setupAuthMocks(t)
		defer ctrl.Finish()

		user := validUserForAuth(t)
		saveErr := errors.New("session save error")
		email, _ := vo.NewEmail("test@example.com")

		setupTxManagerExecute(txManager)
		userRepo.EXPECT().
			FindByEmail(gomock.Any(), gomock.Eq(email)).
			Return(user, nil)
		sessionRepo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(saveErr)

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		testEmail, _ := vo.NewEmail("test@example.com")
		testPassword, _ := vo.NewPassword("password123")
		_, err := uc.Login(context.Background(), testEmail, testPassword)

		if !errors.Is(err, saveErr) {
			t.Errorf("got %v, want saveErr", err)
		}
	})
}

// =============================================================================
// Logout テスト
// =============================================================================

// TestAuthUsecase_Logout はログアウト機能のテスト
func TestAuthUsecase_Logout(t *testing.T) {
	t.Run("正常系_ログアウト成功", func(t *testing.T) {
		userRepo, sessionRepo, txManager, ctrl := setupAuthMocks(t)
		defer ctrl.Finish()

		sid := validSessionID(t)

		setupTxManagerExecute(txManager)
		sessionRepo.EXPECT().
			DeleteByID(gomock.Any(), gomock.Eq(sid)).
			Return(nil)

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		err := uc.Logout(context.Background(), sid)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("異常系_セッション削除エラー", func(t *testing.T) {
		userRepo, sessionRepo, txManager, ctrl := setupAuthMocks(t)
		defer ctrl.Finish()

		sid := validSessionID(t)
		deleteErr := errors.New("delete error")

		setupTxManagerExecute(txManager)
		sessionRepo.EXPECT().
			DeleteByID(gomock.Any(), gomock.Eq(sid)).
			Return(deleteErr)

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		err := uc.Logout(context.Background(), sid)

		if !errors.Is(err, deleteErr) {
			t.Errorf("got %v, want deleteErr", err)
		}
	})
}

// =============================================================================
// ValidateSession テスト
// =============================================================================

// TestAuthUsecase_ValidateSession はセッション検証機能のテスト
func TestAuthUsecase_ValidateSession(t *testing.T) {
	t.Run("正常系_セッション検証成功", func(t *testing.T) {
		userRepo, sessionRepo, txManager, ctrl := setupAuthMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		session := validSession(t, userID)

		sessionRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(session.ID())).
			Return(session, nil)

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		result, err := uc.ValidateSession(context.Background(), session.ID())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("result should not be nil")
		}
		if result.ID().String() != session.ID().String() {
			t.Errorf("session id mismatch: got %v, want %v", result.ID().String(), session.ID().String())
		}
	})

	t.Run("異常系_セッションが見つからない", func(t *testing.T) {
		userRepo, sessionRepo, txManager, ctrl := setupAuthMocks(t)
		defer ctrl.Finish()

		sid := validSessionID(t)

		sessionRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(sid)).
			Return(nil, nil)

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.ValidateSession(context.Background(), sid)

		if !errors.Is(err, domainErrors.ErrSessionNotFound) {
			t.Errorf("got %v, want ErrSessionNotFound", err)
		}
	})

	t.Run("異常系_セッション有効期限切れ", func(t *testing.T) {
		userRepo, sessionRepo, txManager, ctrl := setupAuthMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		// 期限切れのセッションを作成
		expiredSession, err := entity.ReconstructSession(
			validSessionID(t).String(),
			userID.String(),
			time.Now().AddDate(0, 0, -1), // 1日前に期限切れ
			time.Now().AddDate(0, 0, -8), // 8日前に作成
		)
		if err != nil {
			t.Fatalf("failed to create expired session: %v", err)
		}

		sessionRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(expiredSession.ID())).
			Return(expiredSession, nil)

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err = uc.ValidateSession(context.Background(), expiredSession.ID())

		if !errors.Is(err, domainErrors.ErrSessionExpired) {
			t.Errorf("got %v, want ErrSessionExpired", err)
		}
	})

	t.Run("異常系_リポジトリエラー", func(t *testing.T) {
		userRepo, sessionRepo, txManager, ctrl := setupAuthMocks(t)
		defer ctrl.Finish()

		sid := validSessionID(t)
		repoErr := errors.New("db error")

		sessionRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(sid)).
			Return(nil, repoErr)

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.ValidateSession(context.Background(), sid)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})
}
