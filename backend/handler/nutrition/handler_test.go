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
)

func init() {
	gin.SetMode(gin.TestMode)
}

// MockNutritionUsecase はNutritionUsecaseのモック実装
type MockNutritionUsecase struct {
	GetTodayPfcFunc func(ctx context.Context, userID vo.UserID) (*usecase.TodayPfcOutput, error)
}

func (m *MockNutritionUsecase) GetTodayPfc(ctx context.Context, userID vo.UserID) (*usecase.TodayPfcOutput, error) {
	if m.GetTodayPfcFunc != nil {
		return m.GetTodayPfcFunc(ctx, userID)
	}
	return nil, nil
}

func TestNutritionHandler_GetTodayPfc(t *testing.T) {
	t.Run("正常系_今日のPFC取得成功", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"
		now := time.Now()

		mockUsecase := &MockNutritionUsecase{
			GetTodayPfcFunc: func(ctx context.Context, userID vo.UserID) (*usecase.TodayPfcOutput, error) {
				currentPfc := vo.NewPfc(50, 30, 200)
				targetPfc := vo.NewPfc(100, 50, 300)
				return &usecase.TodayPfcOutput{
					Date:       now,
					CurrentPfc: currentPfc,
					TargetPfc:  targetPfc,
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

		// 日付検証
		expectedDate := now.Format("2006-01-02")
		if resp.Date != expectedDate {
			t.Errorf("date = %s, want %s", resp.Date, expectedDate)
		}

		// CurrentPfc検証
		if resp.CurrentPfc.Protein != 50 {
			t.Errorf("current protein = %d, want 50", resp.CurrentPfc.Protein)
		}
		if resp.CurrentPfc.Fat != 30 {
			t.Errorf("current fat = %d, want 30", resp.CurrentPfc.Fat)
		}
		if resp.CurrentPfc.Carbohydrates != 200 {
			t.Errorf("current carbs = %d, want 200", resp.CurrentPfc.Carbohydrates)
		}

		// TargetPfc検証
		if resp.TargetPfc.Protein != 100 {
			t.Errorf("target protein = %d, want 100", resp.TargetPfc.Protein)
		}
		if resp.TargetPfc.Fat != 50 {
			t.Errorf("target fat = %d, want 50", resp.TargetPfc.Fat)
		}
		if resp.TargetPfc.Carbohydrates != 300 {
			t.Errorf("target carbs = %d, want 300", resp.TargetPfc.Carbohydrates)
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
