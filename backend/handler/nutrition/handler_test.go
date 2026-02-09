package nutrition_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/common"
	"caltrack/handler/nutrition"
	"caltrack/handler/nutrition/dto"
	"caltrack/usecase/service"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// MockNutritionUsecase はNutritionUsecaseのモック実装
type MockNutritionUsecase struct {
	GetAdviceFunc func(ctx context.Context, userID vo.UserID) (*service.NutritionAdviceOutput, error)
}

func (m *MockNutritionUsecase) GetAdvice(ctx context.Context, userID vo.UserID) (*service.NutritionAdviceOutput, error) {
	if m.GetAdviceFunc != nil {
		return m.GetAdviceFunc(ctx, userID)
	}
	return nil, nil
}

func TestNutritionHandler_GetAdvice(t *testing.T) {
	t.Run("正常系_アドバイス取得成功", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		mockUsecase := &MockNutritionUsecase{
			GetAdviceFunc: func(ctx context.Context, userID vo.UserID) (*service.NutritionAdviceOutput, error) {
				return &service.NutritionAdviceOutput{
					Advice: "バランスの良い食事を心がけましょう。",
				}, nil
			},
		}

		handler := nutrition.NewNutritionHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/nutrition/advice", nil)
		c.Set("userID", userIDStr)

		handler.GetAdvice(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
		}

		var resp dto.AdviceResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Advice == "" {
			t.Error("advice should not be empty")
		}

		expectedAdvice := "バランスの良い食事を心がけましょう。"
		if resp.Advice != expectedAdvice {
			t.Errorf("advice = %s, want %s", resp.Advice, expectedAdvice)
		}
	})

	t.Run("異常系_未認証（userIDがない）", func(t *testing.T) {
		mockUsecase := &MockNutritionUsecase{}
		handler := nutrition.NewNutritionHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/nutrition/advice", nil)
		// userIDを設定しない

		handler.GetAdvice(c)

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

	t.Run("異常系_ユーザーが見つからない", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		mockUsecase := &MockNutritionUsecase{
			GetAdviceFunc: func(ctx context.Context, userID vo.UserID) (*service.NutritionAdviceOutput, error) {
				return nil, domainErrors.ErrUserNotFound
			},
		}

		handler := nutrition.NewNutritionHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/nutrition/advice", nil)
		c.Set("userID", userIDStr)

		handler.GetAdvice(c)

		if w.Code != http.StatusNotFound {
			t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeNotFound {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeNotFound)
		}
	})

	t.Run("異常系_Usecaseエラー", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		mockUsecase := &MockNutritionUsecase{
			GetAdviceFunc: func(ctx context.Context, userID vo.UserID) (*service.NutritionAdviceOutput, error) {
				return nil, errors.New("database connection error")
			},
		}

		handler := nutrition.NewNutritionHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/nutrition/advice", nil)
		c.Set("userID", userIDStr)

		handler.GetAdvice(c)

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
