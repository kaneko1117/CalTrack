package nutrition_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/common"
	"caltrack/handler/nutrition"
	"caltrack/handler/nutrition/dto"
	"caltrack/usecase"
	"caltrack/usecase/service"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// MockNutritionUsecase はNutritionUsecaseのモック実装
type MockNutritionUsecase struct {
	GetAdviceFunc   func(ctx context.Context, userID vo.UserID) (*service.NutritionAdviceOutput, error)
	GetTodayPfcFunc func(ctx context.Context, userID vo.UserID) (*usecase.TodayPfcOutput, error)
}

func (m *MockNutritionUsecase) GetAdvice(ctx context.Context, userID vo.UserID) (*service.NutritionAdviceOutput, error) {
	if m.GetAdviceFunc != nil {
		return m.GetAdviceFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockNutritionUsecase) GetTodayPfc(ctx context.Context, userID vo.UserID) (*usecase.TodayPfcOutput, error) {
	if m.GetTodayPfcFunc != nil {
		return m.GetTodayPfcFunc(ctx, userID)
	}
	return nil, nil
}

func TestNutritionHandler_GetAdvice(t *testing.T) {
	t.Run("正常系_アドバイス取得成功", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		mockUsecase := &MockNutritionUsecase{
			GetAdviceFunc: func(ctx context.Context, userID vo.UserID) (*service.NutritionAdviceOutput, error) {
				return &service.NutritionAdviceOutput{
					Advice: "今日のPFCバランスは良好です。",
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

		if resp.Advice != "今日のPFCバランスは良好です。" {
			t.Errorf("advice = %s, want 今日のPFCバランスは良好です。", resp.Advice)
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

func TestNutritionHandler_GetTodayPfc(t *testing.T) {
	t.Run("正常系_今日のPFC取得成功", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		mockUsecase := &MockNutritionUsecase{
			GetTodayPfcFunc: func(ctx context.Context, userID vo.UserID) (*usecase.TodayPfcOutput, error) {
				return &usecase.TodayPfcOutput{
					Date:       testTime(),
					CurrentPfc: vo.NewPfc(50.0, 30.0, 150.0),
					TargetPfc:  vo.NewPfc(100.0, 50.0, 200.0),
				}, nil
			},
		}

		handler := nutrition.NewNutritionHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/nutrition/today-pfc", nil)
		c.Set("userID", userIDStr)

		handler.GetTodayPfc(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
		}

		var resp dto.TodayPfcResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Current.Protein != 50.0 {
			t.Errorf("current.protein = %f, want 50.0", resp.Current.Protein)
		}
		if resp.Target.Protein != 100.0 {
			t.Errorf("target.protein = %f, want 100.0", resp.Target.Protein)
		}
	})

	t.Run("異常系_未認証（userIDがない）", func(t *testing.T) {
		mockUsecase := &MockNutritionUsecase{}
		handler := nutrition.NewNutritionHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/nutrition/today-pfc", nil)
		// userIDを設定しない

		handler.GetTodayPfc(c)

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
			GetTodayPfcFunc: func(ctx context.Context, userID vo.UserID) (*usecase.TodayPfcOutput, error) {
				return nil, domainErrors.ErrUserNotFound
			},
		}

		handler := nutrition.NewNutritionHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/nutrition/today-pfc", nil)
		c.Set("userID", userIDStr)

		handler.GetTodayPfc(c)

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
			GetTodayPfcFunc: func(ctx context.Context, userID vo.UserID) (*usecase.TodayPfcOutput, error) {
				return nil, errors.New("database connection error")
			},
		}

		handler := nutrition.NewNutritionHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/nutrition/today-pfc", nil)
		c.Set("userID", userIDStr)

		handler.GetTodayPfc(c)

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

// testTime はテスト用の固定時刻を返す
func testTime() time.Time {
	return time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
}
