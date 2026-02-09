package user_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/user"
	"caltrack/usecase"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// MockUserUsecase はUserUsecaseのモック実装
type MockUserUsecase struct {
	RegisterFunc      func(ctx context.Context, user *entity.User) (*entity.User, error)
	GetProfileFunc    func(ctx context.Context, userID vo.UserID) (*entity.User, error)
	UpdateProfileFunc func(ctx context.Context, userID vo.UserID, input usecase.UpdateProfileInput) (*entity.User, error)
}

func (m *MockUserUsecase) Register(ctx context.Context, user *entity.User) (*entity.User, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, user)
	}
	return nil, nil
}

func (m *MockUserUsecase) GetProfile(ctx context.Context, userID vo.UserID) (*entity.User, error) {
	if m.GetProfileFunc != nil {
		return m.GetProfileFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockUserUsecase) UpdateProfile(ctx context.Context, userID vo.UserID, input usecase.UpdateProfileInput) (*entity.User, error) {
	if m.UpdateProfileFunc != nil {
		return m.UpdateProfileFunc(ctx, userID, input)
	}
	return nil, nil
}

func TestUserHandler_Register(t *testing.T) {
	t.Run("正常系_登録成功", func(t *testing.T) {
		testUser := createTestUser()
		mockUC := &MockUserUsecase{
			RegisterFunc: func(ctx context.Context, user *entity.User) (*entity.User, error) {
				return testUser, nil
			},
		}
		handler := user.NewUserHandler(mockUC)

		reqBody := `{
			"email": "test@example.com",
			"password": "password123",
			"nickname": "testuser",
			"weight": 70.5,
			"height": 175.0,
			"birthDate": "1990-01-01",
			"gender": "male",
			"activityLevel": "moderate"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		if w.Code != http.StatusCreated {
			t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
		}
	})

	t.Run("異常系_メールアドレス重複", func(t *testing.T) {
		mockUC := &MockUserUsecase{
			RegisterFunc: func(ctx context.Context, user *entity.User) (*entity.User, error) {
				return nil, domainErrors.ErrEmailAlreadyExists
			},
		}
		handler := user.NewUserHandler(mockUC)

		reqBody := `{
			"email": "test@example.com",
			"password": "password123",
			"nickname": "testuser",
			"weight": 70.5,
			"height": 175.0,
			"birthDate": "1990-01-01",
			"gender": "male",
			"activityLevel": "moderate"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		if w.Code != http.StatusConflict {
			t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
		}
	})

	t.Run("異常系_バリデーションエラー", func(t *testing.T) {
		mockUC := &MockUserUsecase{}
		handler := user.NewUserHandler(mockUC)

		reqBody := `{
			"email": "invalid-email",
			"password": "password123",
			"nickname": "testuser",
			"weight": 70.5,
			"height": 175.0,
			"birthDate": "1990-01-01",
			"gender": "male",
			"activityLevel": "moderate"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("異常系_不正な生年月日フォーマット", func(t *testing.T) {
		mockUC := &MockUserUsecase{}
		handler := user.NewUserHandler(mockUC)

		reqBody := `{
			"email": "test@example.com",
			"password": "password123",
			"nickname": "testuser",
			"weight": 70.5,
			"height": 175.0,
			"birthDate": "invalid-date",
			"gender": "male",
			"activityLevel": "moderate"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/users", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")

		handler.Register(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})
}

// createTestUser はテスト用のユーザーを生成するヘルパー関数
func createTestUser() *entity.User {
	user, err := entity.ReconstructUser(
		"550e8400-e29b-41d4-a716-446655440000",
		"test@example.com",
		"$2a$10$hashedpassword",
		"TestUser",
		70.0,
		170.0,
		time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		"male",
		"moderate",
		time.Now(),
		time.Now(),
	)
	if err != nil {
		panic(err)
	}
	return user
}

func TestUserHandler_UpdateProfile(t *testing.T) {
	t.Run("正常系_プロフィール更新成功", func(t *testing.T) {
		testUser := createTestUser()
		mockUC := &MockUserUsecase{
			UpdateProfileFunc: func(ctx context.Context, userID vo.UserID, input usecase.UpdateProfileInput) (*entity.User, error) {
				// 更新後のユーザーを返す
				testUser.UpdateProfile("UpdatedNickname", 175.0, 72.5, "active")
				return testUser, nil
			},
		}
		handler := user.NewUserHandler(mockUC)

		reqBody := `{
			"nickname": "UpdatedNickname",
			"height": 175.0,
			"weight": 72.5,
			"activityLevel": "active"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/users/profile", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", testUser.ID().String())

		handler.UpdateProfile(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		if response["userId"] != testUser.ID().String() {
			t.Errorf("userId = %v, want %v", response["userId"], testUser.ID().String())
		}
		if response["nickname"] != "UpdatedNickname" {
			t.Errorf("nickname = %v, want UpdatedNickname", response["nickname"])
		}
		if response["height"] != 175.0 {
			t.Errorf("height = %v, want 175.0", response["height"])
		}
		if response["weight"] != 72.5 {
			t.Errorf("weight = %v, want 72.5", response["weight"])
		}
		if response["activityLevel"] != "active" {
			t.Errorf("activityLevel = %v, want active", response["activityLevel"])
		}
	})

	t.Run("異常系_認証なし", func(t *testing.T) {
		mockUC := &MockUserUsecase{}
		handler := user.NewUserHandler(mockUC)

		reqBody := `{
			"nickname": "UpdatedNickname",
			"height": 175.0,
			"weight": 72.5,
			"activityLevel": "active"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/users/profile", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		// userIDをセットしない

		handler.UpdateProfile(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("異常系_無効なリクエストボディ", func(t *testing.T) {
		testUser := createTestUser()
		mockUC := &MockUserUsecase{}
		handler := user.NewUserHandler(mockUC)

		reqBody := `invalid json`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/users/profile", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", testUser.ID().String())

		handler.UpdateProfile(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("異常系_ユーザーが見つからない", func(t *testing.T) {
		testUser := createTestUser()
		mockUC := &MockUserUsecase{
			UpdateProfileFunc: func(ctx context.Context, userID vo.UserID, input usecase.UpdateProfileInput) (*entity.User, error) {
				return nil, domainErrors.ErrUserNotFound
			},
		}
		handler := user.NewUserHandler(mockUC)

		reqBody := `{
			"nickname": "UpdatedNickname",
			"height": 175.0,
			"weight": 72.5,
			"activityLevel": "active"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/users/profile", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", testUser.ID().String())

		handler.UpdateProfile(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("異常系_バリデーションエラー_ニックネーム空", func(t *testing.T) {
		testUser := createTestUser()
		mockUC := &MockUserUsecase{
			UpdateProfileFunc: func(ctx context.Context, userID vo.UserID, input usecase.UpdateProfileInput) (*entity.User, error) {
				return nil, domainErrors.ErrNicknameRequired
			},
		}
		handler := user.NewUserHandler(mockUC)

		reqBody := `{
			"nickname": "",
			"height": 175.0,
			"weight": 72.5,
			"activityLevel": "active"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/users/profile", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", testUser.ID().String())

		handler.UpdateProfile(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("異常系_バリデーションエラー_不正な活動レベル", func(t *testing.T) {
		testUser := createTestUser()
		mockUC := &MockUserUsecase{
			UpdateProfileFunc: func(ctx context.Context, userID vo.UserID, input usecase.UpdateProfileInput) (*entity.User, error) {
				return nil, domainErrors.ErrInvalidActivityLevel
			},
		}
		handler := user.NewUserHandler(mockUC)

		reqBody := `{
			"nickname": "UpdatedNickname",
			"height": 175.0,
			"weight": 72.5,
			"activityLevel": "invalid"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/users/profile", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", testUser.ID().String())

		handler.UpdateProfile(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}
	})

	t.Run("異常系_DB更新失敗", func(t *testing.T) {
		testUser := createTestUser()
		mockUC := &MockUserUsecase{
			UpdateProfileFunc: func(ctx context.Context, userID vo.UserID, input usecase.UpdateProfileInput) (*entity.User, error) {
				return nil, errors.New("database error")
			},
		}
		handler := user.NewUserHandler(mockUC)

		reqBody := `{
			"nickname": "UpdatedNickname",
			"height": 175.0,
			"weight": 72.5,
			"activityLevel": "active"
		}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPatch, "/api/v1/users/profile", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", testUser.ID().String())

		handler.UpdateProfile(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}

func TestUserHandler_GetProfile(t *testing.T) {
	t.Run("正常系_プロフィール取得成功", func(t *testing.T) {
		testUser := createTestUser()
		mockUC := &MockUserUsecase{
			GetProfileFunc: func(ctx context.Context, userID vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		handler := user.NewUserHandler(mockUC)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		c.Set("userID", testUser.ID().String())

		handler.GetProfile(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d, body: %s", w.Code, http.StatusOK, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		// 全7フィールド検証
		if response["email"] != testUser.Email().String() {
			t.Errorf("email = %v, want %v", response["email"], testUser.Email().String())
		}
		if response["nickname"] != testUser.Nickname().String() {
			t.Errorf("nickname = %v, want %v", response["nickname"], testUser.Nickname().String())
		}
		if response["weight"] != 70.0 {
			t.Errorf("weight = %v, want 70.0", response["weight"])
		}
		if response["height"] != 170.0 {
			t.Errorf("height = %v, want 170.0", response["height"])
		}
		if response["birthDate"] != "1990-01-01" {
			t.Errorf("birthDate = %v, want 1990-01-01", response["birthDate"])
		}
		if response["gender"] != "male" {
			t.Errorf("gender = %v, want male", response["gender"])
		}
		if response["activityLevel"] != "moderate" {
			t.Errorf("activityLevel = %v, want moderate", response["activityLevel"])
		}

		// userIdが含まれないこと確認
		if _, exists := response["userId"]; exists {
			t.Errorf("userId should not be included in response")
		}

		// createdAt/updatedAtが含まれないこと確認
		if _, exists := response["createdAt"]; exists {
			t.Errorf("createdAt should not be included in response")
		}
		if _, exists := response["updatedAt"]; exists {
			t.Errorf("updatedAt should not be included in response")
		}
	})

	t.Run("異常系_認証なし", func(t *testing.T) {
		mockUC := &MockUserUsecase{}
		handler := user.NewUserHandler(mockUC)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		// userIDをセットしない

		handler.GetProfile(c)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
		}
	})

	t.Run("異常系_ユーザーが見つからない", func(t *testing.T) {
		testUser := createTestUser()
		mockUC := &MockUserUsecase{
			GetProfileFunc: func(ctx context.Context, userID vo.UserID) (*entity.User, error) {
				return nil, domainErrors.ErrUserNotFound
			},
		}
		handler := user.NewUserHandler(mockUC)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		c.Set("userID", testUser.ID().String())

		handler.GetProfile(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("異常系_DB取得失敗", func(t *testing.T) {
		testUser := createTestUser()
		mockUC := &MockUserUsecase{
			GetProfileFunc: func(ctx context.Context, userID vo.UserID) (*entity.User, error) {
				return nil, errors.New("database error")
			},
		}
		handler := user.NewUserHandler(mockUC)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/users/profile", nil)
		c.Set("userID", testUser.ID().String())

		handler.GetProfile(c)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
		}
	})
}
