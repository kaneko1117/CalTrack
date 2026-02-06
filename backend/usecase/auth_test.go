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

// mockSessionRepository はSessionRepositoryのモック実装
type mockSessionRepository struct {
	save           func(ctx context.Context, session *entity.Session) error
	findByID       func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error)
	deleteByID     func(ctx context.Context, sessionID vo.SessionID) error
	deleteByUserID func(ctx context.Context, userID vo.UserID) error
}

func (m *mockSessionRepository) Save(ctx context.Context, session *entity.Session) error {
	return m.save(ctx, session)
}

func (m *mockSessionRepository) FindByID(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
	return m.findByID(ctx, sessionID)
}

func (m *mockSessionRepository) DeleteByID(ctx context.Context, sessionID vo.SessionID) error {
	return m.deleteByID(ctx, sessionID)
}

func (m *mockSessionRepository) DeleteByUserID(ctx context.Context, userID vo.UserID) error {
	return m.deleteByUserID(ctx, userID)
}

// mockUserRepositoryForAuth はUserRepositoryのモック実装（auth用）
type mockUserRepositoryForAuth struct {
	findByEmail   func(ctx context.Context, email vo.Email) (*entity.User, error)
	existsByEmail func(ctx context.Context, email vo.Email) (bool, error)
	save          func(ctx context.Context, user *entity.User) error
	findByID      func(ctx context.Context, id vo.UserID) (*entity.User, error)
}

func (m *mockUserRepositoryForAuth) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	return m.findByEmail(ctx, email)
}

func (m *mockUserRepositoryForAuth) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	return m.existsByEmail(ctx, email)
}

func (m *mockUserRepositoryForAuth) Save(ctx context.Context, user *entity.User) error {
	return m.save(ctx, user)
}

func (m *mockUserRepositoryForAuth) FindByID(ctx context.Context, id vo.UserID) (*entity.User, error) {
	if m.findByID != nil {
		return m.findByID(ctx, id)
	}
	return nil, nil
}

func (m *mockUserRepositoryForAuth) Update(ctx context.Context, user *entity.User) error {
	return nil
}

// mockTxManager はTransactionManagerのモック実装
type mockTxManager struct{}

func (m *mockTxManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
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
		user := validUserForAuth(t)

		userRepo := &mockUserRepositoryForAuth{
			findByEmail: func(ctx context.Context, email vo.Email) (*entity.User, error) {
				return user, nil
			},
		}
		sessionRepo := &mockSessionRepository{
			save: func(ctx context.Context, session *entity.Session) error {
				return nil
			},
		}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		output, err := uc.Login(context.Background(), usecase.LoginInput{
			Email:    "test@example.com",
			Password: "password123",
		})

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

	t.Run("異常系_無効なメールアドレス形式", func(t *testing.T) {
		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.Login(context.Background(), usecase.LoginInput{
			Email:    "invalid-email",
			Password: "password123",
		})

		if !errors.Is(err, domainErrors.ErrInvalidCredentials) {
			t.Errorf("got %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("異常系_空のメールアドレス", func(t *testing.T) {
		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.Login(context.Background(), usecase.LoginInput{
			Email:    "",
			Password: "password123",
		})

		if !errors.Is(err, domainErrors.ErrInvalidCredentials) {
			t.Errorf("got %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("異常系_無効なパスワード形式", func(t *testing.T) {
		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.Login(context.Background(), usecase.LoginInput{
			Email:    "test@example.com",
			Password: "short",
		})

		if !errors.Is(err, domainErrors.ErrInvalidCredentials) {
			t.Errorf("got %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("異常系_空のパスワード", func(t *testing.T) {
		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.Login(context.Background(), usecase.LoginInput{
			Email:    "test@example.com",
			Password: "",
		})

		if !errors.Is(err, domainErrors.ErrInvalidCredentials) {
			t.Errorf("got %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("異常系_ユーザーが見つからない", func(t *testing.T) {
		userRepo := &mockUserRepositoryForAuth{
			findByEmail: func(ctx context.Context, email vo.Email) (*entity.User, error) {
				return nil, nil
			},
		}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.Login(context.Background(), usecase.LoginInput{
			Email:    "notfound@example.com",
			Password: "password123",
		})

		if !errors.Is(err, domainErrors.ErrInvalidCredentials) {
			t.Errorf("got %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("異常系_パスワードが一致しない", func(t *testing.T) {
		user := validUserForAuth(t)

		userRepo := &mockUserRepositoryForAuth{
			findByEmail: func(ctx context.Context, email vo.Email) (*entity.User, error) {
				return user, nil
			},
		}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.Login(context.Background(), usecase.LoginInput{
			Email:    "test@example.com",
			Password: "wrongpassword",
		})

		if !errors.Is(err, domainErrors.ErrInvalidCredentials) {
			t.Errorf("got %v, want ErrInvalidCredentials", err)
		}
	})

	t.Run("異常系_ユーザーリポジトリエラー", func(t *testing.T) {
		repoErr := errors.New("db error")
		userRepo := &mockUserRepositoryForAuth{
			findByEmail: func(ctx context.Context, email vo.Email) (*entity.User, error) {
				return nil, repoErr
			},
		}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.Login(context.Background(), usecase.LoginInput{
			Email:    "test@example.com",
			Password: "password123",
		})

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_セッション保存エラー", func(t *testing.T) {
		user := validUserForAuth(t)
		saveErr := errors.New("session save error")

		userRepo := &mockUserRepositoryForAuth{
			findByEmail: func(ctx context.Context, email vo.Email) (*entity.User, error) {
				return user, nil
			},
		}
		sessionRepo := &mockSessionRepository{
			save: func(ctx context.Context, session *entity.Session) error {
				return saveErr
			},
		}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.Login(context.Background(), usecase.LoginInput{
			Email:    "test@example.com",
			Password: "password123",
		})

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
		sid := validSessionID(t)

		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{
			deleteByID: func(ctx context.Context, sessionID vo.SessionID) error {
				return nil
			},
		}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		err := uc.Logout(context.Background(), sid.String())

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("異常系_無効なセッションID", func(t *testing.T) {
		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		err := uc.Logout(context.Background(), "invalid-session-id")

		if !errors.Is(err, domainErrors.ErrInvalidSessionID) {
			t.Errorf("got %v, want ErrInvalidSessionID", err)
		}
	})

	t.Run("異常系_空のセッションID", func(t *testing.T) {
		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		err := uc.Logout(context.Background(), "")

		if !errors.Is(err, domainErrors.ErrInvalidSessionID) {
			t.Errorf("got %v, want ErrInvalidSessionID", err)
		}
	})

	t.Run("異常系_セッション削除エラー", func(t *testing.T) {
		sid := validSessionID(t)
		deleteErr := errors.New("delete error")

		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{
			deleteByID: func(ctx context.Context, sessionID vo.SessionID) error {
				return deleteErr
			},
		}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		err := uc.Logout(context.Background(), sid.String())

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
		userID := vo.NewUserID()
		session := validSession(t, userID)

		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{
			findByID: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return session, nil
			},
		}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		result, err := uc.ValidateSession(context.Background(), session.ID().String())

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

	t.Run("異常系_無効なセッションID", func(t *testing.T) {
		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.ValidateSession(context.Background(), "invalid-session-id")

		if !errors.Is(err, domainErrors.ErrInvalidSessionID) {
			t.Errorf("got %v, want ErrInvalidSessionID", err)
		}
	})

	t.Run("異常系_空のセッションID", func(t *testing.T) {
		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.ValidateSession(context.Background(), "")

		if !errors.Is(err, domainErrors.ErrInvalidSessionID) {
			t.Errorf("got %v, want ErrInvalidSessionID", err)
		}
	})

	t.Run("異常系_セッションが見つからない", func(t *testing.T) {
		sid := validSessionID(t)

		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{
			findByID: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return nil, nil
			},
		}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.ValidateSession(context.Background(), sid.String())

		if !errors.Is(err, domainErrors.ErrSessionNotFound) {
			t.Errorf("got %v, want ErrSessionNotFound", err)
		}
	})

	t.Run("異常系_セッション有効期限切れ", func(t *testing.T) {
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

		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{
			findByID: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return expiredSession, nil
			},
		}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err = uc.ValidateSession(context.Background(), expiredSession.ID().String())

		if !errors.Is(err, domainErrors.ErrSessionExpired) {
			t.Errorf("got %v, want ErrSessionExpired", err)
		}
	})

	t.Run("異常系_リポジトリエラー", func(t *testing.T) {
		sid := validSessionID(t)
		repoErr := errors.New("db error")

		userRepo := &mockUserRepositoryForAuth{}
		sessionRepo := &mockSessionRepository{
			findByID: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return nil, repoErr
			},
		}
		txManager := &mockTxManager{}

		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		_, err := uc.ValidateSession(context.Background(), sid.String())

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})
}
