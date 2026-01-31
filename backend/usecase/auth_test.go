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

// TestLogin_Success はログイン成功のテスト
func TestLogin_Success(t *testing.T) {
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
}

// TestLogin_InvalidEmailFormat は無効なメールアドレス形式のテスト
func TestLogin_InvalidEmailFormat(t *testing.T) {
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
}

// TestLogin_EmptyEmail は空のメールアドレスのテスト
func TestLogin_EmptyEmail(t *testing.T) {
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
}

// TestLogin_InvalidPasswordFormat は無効なパスワード形式のテスト（短すぎる）
func TestLogin_InvalidPasswordFormat(t *testing.T) {
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
}

// TestLogin_EmptyPassword は空のパスワードのテスト
func TestLogin_EmptyPassword(t *testing.T) {
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
}

// TestLogin_UserNotFound はユーザーが見つからない場合のテスト
func TestLogin_UserNotFound(t *testing.T) {
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
}

// TestLogin_PasswordMismatch はパスワードが一致しない場合のテスト
func TestLogin_PasswordMismatch(t *testing.T) {
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
}

// TestLogin_UserRepositoryError はユーザーリポジトリエラーのテスト
func TestLogin_UserRepositoryError(t *testing.T) {
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
}

// TestLogin_SessionSaveError はセッション保存エラーのテスト
func TestLogin_SessionSaveError(t *testing.T) {
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
}

// =============================================================================
// Logout テスト
// =============================================================================

// TestLogout_Success はログアウト成功のテスト
func TestLogout_Success(t *testing.T) {
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
}

// TestLogout_InvalidSessionID は無効なセッションIDのテスト
func TestLogout_InvalidSessionID(t *testing.T) {
	userRepo := &mockUserRepositoryForAuth{}
	sessionRepo := &mockSessionRepository{}
	txManager := &mockTxManager{}

	uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
	err := uc.Logout(context.Background(), "invalid-session-id")

	if !errors.Is(err, domainErrors.ErrInvalidSessionID) {
		t.Errorf("got %v, want ErrInvalidSessionID", err)
	}
}

// TestLogout_EmptySessionID は空のセッションIDのテスト
func TestLogout_EmptySessionID(t *testing.T) {
	userRepo := &mockUserRepositoryForAuth{}
	sessionRepo := &mockSessionRepository{}
	txManager := &mockTxManager{}

	uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
	err := uc.Logout(context.Background(), "")

	if !errors.Is(err, domainErrors.ErrInvalidSessionID) {
		t.Errorf("got %v, want ErrInvalidSessionID", err)
	}
}

// TestLogout_DeleteError はセッション削除エラーのテスト
func TestLogout_DeleteError(t *testing.T) {
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
}

// =============================================================================
// ValidateSession テスト
// =============================================================================

// TestValidateSession_Success はセッション検証成功のテスト
func TestValidateSession_Success(t *testing.T) {
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
}

// TestValidateSession_InvalidSessionID は無効なセッションIDのテスト
func TestValidateSession_InvalidSessionID(t *testing.T) {
	userRepo := &mockUserRepositoryForAuth{}
	sessionRepo := &mockSessionRepository{}
	txManager := &mockTxManager{}

	uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
	_, err := uc.ValidateSession(context.Background(), "invalid-session-id")

	if !errors.Is(err, domainErrors.ErrInvalidSessionID) {
		t.Errorf("got %v, want ErrInvalidSessionID", err)
	}
}

// TestValidateSession_EmptySessionID は空のセッションIDのテスト
func TestValidateSession_EmptySessionID(t *testing.T) {
	userRepo := &mockUserRepositoryForAuth{}
	sessionRepo := &mockSessionRepository{}
	txManager := &mockTxManager{}

	uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
	_, err := uc.ValidateSession(context.Background(), "")

	if !errors.Is(err, domainErrors.ErrInvalidSessionID) {
		t.Errorf("got %v, want ErrInvalidSessionID", err)
	}
}

// TestValidateSession_SessionNotFound はセッションが見つからない場合のテスト
func TestValidateSession_SessionNotFound(t *testing.T) {
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
}

// TestValidateSession_SessionExpired は有効期限切れセッションのテスト
func TestValidateSession_SessionExpired(t *testing.T) {
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
}

// TestValidateSession_RepositoryError はリポジトリエラーのテスト
func TestValidateSession_RepositoryError(t *testing.T) {
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
}
