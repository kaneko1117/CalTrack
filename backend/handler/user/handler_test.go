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

func TestUserHandler_Register(t *testing.T) {
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
		handler := user.NewUserHandler(uc)

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
		handler := user.NewUserHandler(uc)

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
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
			update: func(ctx context.Context, user *entity.User) error {
				return nil
			},
		}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, domainErrors.ErrUserNotFound
			},
		}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
			update: func(ctx context.Context, user *entity.User) error {
				return errors.New("database error")
			},
		}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, domainErrors.ErrUserNotFound
			},
		}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
		repo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, errors.New("database error")
			},
		}
		txManager := &mockTransactionManager{}
		uc := usecase.NewUserUsecase(repo, txManager)
		handler := user.NewUserHandler(uc)

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
