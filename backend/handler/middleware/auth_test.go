package middleware_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/common"
	"caltrack/handler/middleware"
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
	findByID      func(ctx context.Context, id vo.UserID) (*entity.User, error)
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

func (m *mockUserRepository) FindByID(ctx context.Context, id vo.UserID) (*entity.User, error) {
	if m.findByID != nil {
		return m.findByID(ctx, id)
	}
	return nil, nil
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

		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{
			findByID: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return testSession, nil
			},
		}
		txManager := &mockTransactionManager{}
		authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)

		// ミドルウェアを設定したルーターを作成
		r := gin.New()
		r.Use(middleware.AuthMiddleware(authUsecase))
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
		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTransactionManager{}
		authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)

		r := gin.New()
		r.Use(middleware.AuthMiddleware(authUsecase))
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
		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{}
		txManager := &mockTransactionManager{}
		authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)

		r := gin.New()
		r.Use(middleware.AuthMiddleware(authUsecase))
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
		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{
			findByID: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return nil, nil // セッションが見つからない
			},
		}
		txManager := &mockTransactionManager{}
		authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)

		r := gin.New()
		r.Use(middleware.AuthMiddleware(authUsecase))
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
		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{
			findByID: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				// 期限切れセッションを返す
				expiredSession, err := entity.ReconstructSession(
					sessionID.String(),
					"550e8400-e29b-41d4-a716-446655440000",
					time.Now().Add(-24*time.Hour), // 過去の時刻（期限切れ）
					time.Now().Add(-48*time.Hour), // 作成日時
				)
				if err != nil {
					return nil, err
				}
				return expiredSession, nil
			},
		}
		txManager := &mockTransactionManager{}
		authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)

		r := gin.New()
		r.Use(middleware.AuthMiddleware(authUsecase))
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
		userRepo := &mockUserRepository{}
		sessionRepo := &mockSessionRepository{
			findByID: func(ctx context.Context, sessionID vo.SessionID) (*entity.Session, error) {
				return nil, domainErrors.ErrSessionIDGenerationFailed
			},
		}
		txManager := &mockTransactionManager{}
		authUsecase := usecase.NewAuthUsecase(userRepo, sessionRepo, txManager)

		r := gin.New()
		r.Use(middleware.AuthMiddleware(authUsecase))
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
