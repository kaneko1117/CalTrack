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

// MockAuthUsecase はAuthUsecaseのモック実装
type MockAuthUsecase struct {
	LoginFunc           func(ctx context.Context, email vo.Email, password vo.Password) (*usecase.LoginOutput, error)
	LogoutFunc          func(ctx context.Context, sessionID vo.SessionID) error
	ValidateSessionFunc func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error)
}

func (m *MockAuthUsecase) Login(ctx context.Context, email vo.Email, password vo.Password) (*usecase.LoginOutput, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, email, password)
	}
	return nil, nil
}

func (m *MockAuthUsecase) Logout(ctx context.Context, sessionID vo.SessionID) error {
	if m.LogoutFunc != nil {
		return m.LogoutFunc(ctx, sessionID)
	}
	return nil
}

func (m *MockAuthUsecase) ValidateSession(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
	if m.ValidateSessionFunc != nil {
		return m.ValidateSessionFunc(ctx, sessionID)
	}
	return nil, nil
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

		mockUC := &MockAuthUsecase{
			LoginFunc: func(ctx context.Context, email vo.Email, password vo.Password) (*usecase.LoginOutput, error) {
				// テスト用のセッションを作成
				session, err := entity.NewSessionWithUserID(testUser.ID())
				if err != nil {
					return nil, err
				}
				return &usecase.LoginOutput{
					User:    testUser,
					Session: session,
				}, nil
			},
		}
		handler := auth.NewAuthHandler(mockUC)

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
		mockUC := &MockAuthUsecase{
			LoginFunc: func(ctx context.Context, email vo.Email, password vo.Password) (*usecase.LoginOutput, error) {
				return nil, domainErrors.ErrInvalidCredentials
			},
		}
		handler := auth.NewAuthHandler(mockUC)

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
		mockUC := &MockAuthUsecase{
			LoginFunc: func(ctx context.Context, email vo.Email, password vo.Password) (*usecase.LoginOutput, error) {
				return nil, domainErrors.ErrInvalidCredentials
			},
		}
		handler := auth.NewAuthHandler(mockUC)

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
		mockUC := &MockAuthUsecase{}
		handler := auth.NewAuthHandler(mockUC)

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
		mockUC := &MockAuthUsecase{}
		handler := auth.NewAuthHandler(mockUC)

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
		mockUC := &MockAuthUsecase{
			LogoutFunc: func(ctx context.Context, sessionID vo.SessionID) error {
				return nil
			},
		}
		handler := auth.NewAuthHandler(mockUC)

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
		mockUC := &MockAuthUsecase{}
		handler := auth.NewAuthHandler(mockUC)

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
		mockUC := &MockAuthUsecase{
			LogoutFunc: func(ctx context.Context, sessionID vo.SessionID) error {
				return domainErrors.ErrInvalidSessionID
			},
		}
		handler := auth.NewAuthHandler(mockUC)

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
		mockUC := &MockAuthUsecase{
			LogoutFunc: func(ctx context.Context, sessionID vo.SessionID) error {
				return domainErrors.ErrSessionNotFound
			},
		}
		handler := auth.NewAuthHandler(mockUC)

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
