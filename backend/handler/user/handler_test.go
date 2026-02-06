package user_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"caltrack/domain/entity"
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
