package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/auth"
	"caltrack/handler/common"
	"caltrack/usecase"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockUserRepository はUserRepositoryのモック実装
type mockUserRepository struct {
	existsByEmail func(ctx context.Context, email vo.Email) (bool, error)
	save          func(ctx context.Context, user *entity.User) error
	findByEmail   func(ctx context.Context, email vo.Email) (*entity.User, error)
}

func (m *mockUserRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	return m.existsByEmail(ctx, email)
}

func (m *mockUserRepository) Save(ctx context.Context, u *entity.User) error {
	return m.save(ctx, u)
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	return m.findByEmail(ctx, email)
}

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

// mockTransactionManager はTransactionManagerのモック実装
type mockTransactionManager struct{}

func (m *mockTransactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// createTestUser はテスト用ユーザーを作成するヘルパー関数
func createTestUser(t *testing.T, email, password string) *entity.User {
	t.Helper()
	// パスワードをハッシュ化するためにNewUserを使用
	user, errs := entity.NewUser(
		email,
		password,
		"testuser",
		70.5,
		175.0,
		time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		"male",
		"moderate",
	)
	if len(errs) > 0 {
		t.Fatalf("failed to create test user: %v", errs)
	}
	return user
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("正常系_ログイン成功", func(t *testing.T) {
		testEmail := "test@example.com"
		testPassword := "password123"
		testUser := createTestUser(t, testEmail, testPassword)

		userRepo := &mockUserRepository{
			findByEmail: func(ctx context.Context, email vo.Email) (*entity.User, error) {
				return testUser, nil
			},
		}
		sessionRepo := &mockSessionRepository{
			save: func(ctx context.Context, session *entity.Session) error {
				return nil
			},
		}
		txManager := &mockTransactionManager{}
		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		handler := auth.NewAuthHandler(uc)

		reqBody := `{"email": "test@example.com", "password": "password123"}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
		}

		// レスポンスボディの検証
		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp["email"] != testEmail {
			t.Errorf("email = %s, want %s", resp["email"], testEmail)
		}
		if resp["nickname"] != "testuser" {
			t.Errorf("nickname = %s, want %s", resp["nickname"], "testuser")
		}
		if resp["userId"] == "" {
			t.Error("userId should not be empty")
		}

		// Cookieの検証
		cookies := w.Result().Cookies()
		var sessionCookie *http.Cookie
		for _, cookie := range cookies {
			if cookie.Name == "session_id" {
				sessionCookie = cookie
				break
			}
		}
		if sessionCookie == nil {
			t.Error("session_id cookie should be set")
		} else {
			if !sessionCookie.HttpOnly {
				t.Error("session_id cookie should be HttpOnly")
			}
			if !sessionCookie.Secure {
				t.Error("session_id cookie should be Secure")
			}
		}
	})

	t.Run("異常系_認証情報不正", func(t *testing.T) {
		testEmail := "test@example.com"
		testPassword := "password123"
		testUser := createTestUser(t, testEmail, testPassword)

		userRepo := &mockUserRepository{
			findByEmail: func(ctx context.Context, email vo.Email) (*entity.User, error) {
				return testUser, nil
			},
		}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTransactionManager{}
		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		handler := auth.NewAuthHandler(uc)

		// 間違ったパスワードでログイン試行
		reqBody := `{"email": "test@example.com", "password": "wrongpassword"}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}

		// エラーレスポンスの検証
		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeInvalidCredentials {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeInvalidCredentials)
		}
	})

	t.Run("異常系_ユーザー未登録", func(t *testing.T) {
		userRepo := &mockUserRepository{
			findByEmail: func(ctx context.Context, email vo.Email) (*entity.User, error) {
				return nil, nil // ユーザーが見つからない
			},
		}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTransactionManager{}
		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		handler := auth.NewAuthHandler(uc)

		reqBody := `{"email": "notfound@example.com", "password": "password123"}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeInvalidCredentials {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeInvalidCredentials)
		}
	})

	t.Run("異常系_メールアドレス形式不正", func(t *testing.T) {
		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTransactionManager{}
		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		handler := auth.NewAuthHandler(uc)

		reqBody := `{"email": "invalid-email", "password": "password123"}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeInvalidCredentials {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeInvalidCredentials)
		}
	})

	t.Run("異常系_リクエストボディ不正", func(t *testing.T) {
		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTransactionManager{}
		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		handler := auth.NewAuthHandler(uc)

		reqBody := `{invalid json}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Login(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeInvalidRequest {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeInvalidRequest)
		}
	})
}

func TestAuthHandler_Logout(t *testing.T) {
	t.Run("正常系_ログアウト成功", func(t *testing.T) {
		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{
			deleteByID: func(ctx context.Context, sessionID vo.SessionID) error {
				return nil
			},
		}
		txManager := &mockTransactionManager{}
		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		handler := auth.NewAuthHandler(uc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)

		// 有効なセッションIDをCookieに設定（44文字のBase64エンコード）
		validSessionID := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
		c.Request.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: validSessionID,
		})

		handler.Logout(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
		}

		// Cookieが削除されていることを確認
		cookies := w.Result().Cookies()
		for _, cookie := range cookies {
			if cookie.Name == "session_id" {
				if cookie.MaxAge >= 0 {
					t.Error("session_id cookie should be deleted (MaxAge < 0)")
				}
				break
			}
		}
	})

	t.Run("正常系_Cookie未設定", func(t *testing.T) {
		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTransactionManager{}
		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		handler := auth.NewAuthHandler(uc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)
		// Cookieを設定しない

		handler.Logout(c)

		// Cookieがなくても成功として扱う
		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}
	})

	t.Run("異常系_セッションID不正", func(t *testing.T) {
		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTransactionManager{}
		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		handler := auth.NewAuthHandler(uc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)

		// 無効なセッションIDをCookieに設定
		c.Request.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "invalid-session-id",
		})

		handler.Logout(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("異常系_セッション削除エラー", func(t *testing.T) {
		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{
			deleteByID: func(ctx context.Context, sessionID vo.SessionID) error {
				return domainErrors.ErrSessionNotFound
			},
		}
		txManager := &mockTransactionManager{}
		uc := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)
		handler := auth.NewAuthHandler(uc)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/auth/logout", nil)

		// 有効なセッションIDをCookieに設定
		validSessionID := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
		c.Request.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: validSessionID,
		})

		handler.Logout(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeUnauthorized {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeUnauthorized)
		}
	})
}
