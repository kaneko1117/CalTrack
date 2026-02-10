package middleware_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/common"
	"caltrack/handler/middleware"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// MockAuthSessionValidator はAuthSessionValidatorのモック実装
type MockAuthSessionValidator struct {
	ValidateSessionFunc func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error)
}

func (m *MockAuthSessionValidator) ValidateSession(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
	if m.ValidateSessionFunc != nil {
		return m.ValidateSessionFunc(ctx, sessionID)
	}
	return nil, nil
}

// createTestSession はテスト用のセッションを作成する
func createTestSession(t *testing.T, userIDStr string) *entity.Session {
	t.Helper()
	userID := vo.ReconstructUserID(userIDStr)
	session, err := entity.NewSessionWithUserID(userID)
	if err != nil {
		t.Fatalf("failed to create test session: %v", err)
	}
	return session
}

func TestAuthMiddleware(t *testing.T) {
	t.Run("正常系_認証成功", func(t *testing.T) {
		testUserID := "550e8400-e29b-41d4-a716-446655440000"
		testSession := createTestSession(t, testUserID)

		mockValidator := &MockAuthSessionValidator{
			ValidateSessionFunc: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return testSession, nil
			},
		}

		// ミドルウェアを設定したルーターを作成
		r := gin.New()
		r.Use(middleware.AuthMiddleware(mockValidator))
		r.GET("/protected", func(c *gin.Context) {
			userID, _ := c.Get("userID")
			c.JSON(http.StatusOK, gin.H{"userID": userID})
		})

		// リクエストを作成
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: testSession.ID().String(),
		})

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
		}

		var resp map[string]string
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp["userID"] != testUserID {
			t.Errorf("userID = %s, want %s", resp["userID"], testUserID)
		}
	})

	t.Run("異常系_Cookieなし", func(t *testing.T) {
		mockValidator := &MockAuthSessionValidator{}

		r := gin.New()
		r.Use(middleware.AuthMiddleware(mockValidator))
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		// Cookieを設定しない

		r.ServeHTTP(w, req)

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

	t.Run("異常系_無効なセッションID", func(t *testing.T) {
		mockValidator := &MockAuthSessionValidator{
			ValidateSessionFunc: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return nil, domainErrors.ErrInvalidSessionID
			},
		}

		r := gin.New()
		r.Use(middleware.AuthMiddleware(mockValidator))
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: "invalid-session-id",
		})

		r.ServeHTTP(w, req)

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

	t.Run("異常系_セッションが見つからない", func(t *testing.T) {
		mockValidator := &MockAuthSessionValidator{
			ValidateSessionFunc: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return nil, domainErrors.ErrSessionNotFound
			},
		}

		r := gin.New()
		r.Use(middleware.AuthMiddleware(mockValidator))
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		// 有効なセッションID形式
		validSessionID := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: validSessionID,
		})

		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("異常系_セッション期限切れ", func(t *testing.T) {
		mockValidator := &MockAuthSessionValidator{
			ValidateSessionFunc: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return nil, domainErrors.ErrSessionExpired
			},
		}

		r := gin.New()
		r.Use(middleware.AuthMiddleware(mockValidator))
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		validSessionID := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: validSessionID,
		})

		r.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeSessionExpired {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeSessionExpired)
		}
	})

	t.Run("異常系_内部エラー", func(t *testing.T) {
		mockValidator := &MockAuthSessionValidator{
			ValidateSessionFunc: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return nil, domainErrors.ErrSessionIDGenerationFailed
			},
		}

		r := gin.New()
		r.Use(middleware.AuthMiddleware(mockValidator))
		r.GET("/protected", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/protected", nil)
		validSessionID := "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
		req.AddCookie(&http.Cookie{
			Name:  "session_id",
			Value: validSessionID,
		})

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeInternalError {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeInternalError)
		}
	})
}
