package record_test

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
	"caltrack/handler/common"
	"caltrack/handler/record"
	"caltrack/handler/record/dto"
	"caltrack/usecase"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// MockRecordUsecase はRecordUsecaseのモック実装
type MockRecordUsecase struct {
	CreateFunc           func(ctx context.Context, record *entity.Record) error
	GetTodayCaloriesFunc func(ctx context.Context, userID vo.UserID) (*usecase.TodayCaloriesOutput, error)
	GetStatisticsFunc    func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) (*usecase.StatisticsOutput, error)
}

func (m *MockRecordUsecase) Create(ctx context.Context, rec *entity.Record) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, rec)
	}
	return nil
}

func (m *MockRecordUsecase) GetTodayCalories(ctx context.Context, userID vo.UserID) (*usecase.TodayCaloriesOutput, error) {
	if m.GetTodayCaloriesFunc != nil {
		return m.GetTodayCaloriesFunc(ctx, userID)
	}
	return nil, nil
}

func (m *MockRecordUsecase) GetStatistics(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) (*usecase.StatisticsOutput, error) {
	if m.GetStatisticsFunc != nil {
		return m.GetStatisticsFunc(ctx, userID, period)
	}
	return nil, nil
}

func TestRecordHandler_Create(t *testing.T) {
	t.Run("正常系_記録が作成される", func(t *testing.T) {
		mockUsecase := &MockRecordUsecase{
			CreateFunc: func(ctx context.Context, rec *entity.Record) error {
				return nil
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		// 現在時刻から少し前の日時を使用（未来日時エラーを回避）
		eatenAt := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
		reqBody := `{
		"eatenAt": "` + eatenAt + `",
		"items": [
			{"name": "ご飯", "calories": 250},
			{"name": "味噌汁", "calories": 50}
		]
	}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/records", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		handler.Create(c)

		if w.Code != http.StatusCreated {
			t.Errorf("status = %d, want %d, body = %s", w.Code, http.StatusCreated, w.Body.String())
		}

		// レスポンスボディの検証
		var resp dto.CreateRecordResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.RecordID == "" {
			t.Error("recordId should not be empty")
		}
		if resp.TotalCalories != 300 {
			t.Errorf("totalCalories = %d, want %d", resp.TotalCalories, 300)
		}
		if len(resp.Items) != 2 {
			t.Errorf("items count = %d, want %d", len(resp.Items), 2)
		}
	})

	t.Run("異常系_認証なし", func(t *testing.T) {
		mockUsecase := &MockRecordUsecase{}
		handler := record.NewRecordHandler(mockUsecase)

		reqBody := `{
		"eatenAt": "2024-01-15T12:00:00Z",
		"items": [{"name": "ご飯", "calories": 250}]
	}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/records", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		// userIDを設定しない

		handler.Create(c)

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

	t.Run("異常系_無効なリクエストボディ", func(t *testing.T) {
		mockUsecase := &MockRecordUsecase{}
		handler := record.NewRecordHandler(mockUsecase)

		reqBody := `{invalid json}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/records", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		handler.Create(c)

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

	t.Run("異常系_明細が空", func(t *testing.T) {
		mockUsecase := &MockRecordUsecase{}
		handler := record.NewRecordHandler(mockUsecase)

		reqBody := `{
		"eatenAt": "2024-01-15T12:00:00Z",
		"items": []
	}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/records", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		handler.Create(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeValidationError {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeValidationError)
		}
		if len(resp.Details) == 0 {
			t.Error("details should not be empty")
		}
	})

	t.Run("異常系_無効なeatenAt形式", func(t *testing.T) {
		mockUsecase := &MockRecordUsecase{}
		handler := record.NewRecordHandler(mockUsecase)

		reqBody := `{
		"eatenAt": "invalid-date-format",
		"items": [{"name": "ご飯", "calories": 250}]
	}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/records", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		handler.Create(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeValidationError {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeValidationError)
		}
	})

	t.Run("異常系_未来日時エラー", func(t *testing.T) {
		mockUsecase := &MockRecordUsecase{}
		handler := record.NewRecordHandler(mockUsecase)

		// 1年後の日時を設定
		futureTime := time.Now().AddDate(1, 0, 0).Format(time.RFC3339)
		reqBody := `{
		"eatenAt": "` + futureTime + `",
		"items": [{"name": "ご飯", "calories": 250}]
	}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/records", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		handler.Create(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeValidationError {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeValidationError)
		}
	})

	t.Run("異常系_カロリーが負の値", func(t *testing.T) {
		mockUsecase := &MockRecordUsecase{}
		handler := record.NewRecordHandler(mockUsecase)

		eatenAt := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
		reqBody := `{
		"eatenAt": "` + eatenAt + `",
		"items": [{"name": "ご飯", "calories": -100}]
	}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/records", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		handler.Create(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeValidationError {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeValidationError)
		}
	})

	t.Run("異常系_食品名が空", func(t *testing.T) {
		mockUsecase := &MockRecordUsecase{}
		handler := record.NewRecordHandler(mockUsecase)

		eatenAt := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
		reqBody := `{
		"eatenAt": "` + eatenAt + `",
		"items": [{"name": "", "calories": 250}]
	}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/records", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		handler.Create(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeValidationError {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeValidationError)
		}
	})

	t.Run("異常系_DB保存失敗", func(t *testing.T) {
		mockUsecase := &MockRecordUsecase{
			CreateFunc: func(ctx context.Context, rec *entity.Record) error {
				return context.DeadlineExceeded
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		eatenAt := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)
		reqBody := `{
		"eatenAt": "` + eatenAt + `",
		"items": [{"name": "ご飯", "calories": 250}]
	}`

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodPost, "/api/v1/records", strings.NewReader(reqBody))
		c.Request.Header.Set("Content-Type", "application/json")
		c.Set("userID", "550e8400-e29b-41d4-a716-446655440000")

		handler.Create(c)

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

func TestRecordHandler_GetToday(t *testing.T) {
	t.Run("正常系_今日のカロリー情報が取得できる", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		// Usecaseのアウトプットを作成
		now := time.Now()
		eatenAt1 := now.Add(-2 * time.Hour)
		eatenAt2 := now.Add(-1 * time.Hour)

		record1 := entity.ReconstructRecord(
			"record-1",
			userIDStr,
			eatenAt1,
			now,
			[]entity.RecordItem{
				*entity.ReconstructRecordItem("item-1", "record-1", "朝食：パン", 300),
				*entity.ReconstructRecordItem("item-2", "record-1", "朝食：コーヒー", 50),
			},
		)
		record2 := entity.ReconstructRecord(
			"record-2",
			userIDStr,
			eatenAt2,
			now,
			[]entity.RecordItem{
				*entity.ReconstructRecordItem("item-3", "record-2", "昼食：ラーメン", 800),
			},
		)

		output := &usecase.TodayCaloriesOutput{
			Date:           now,
			TotalCalories:  1150,
			TargetCalories: 2000,
			Difference:     850,
			Records:        []*entity.Record{record1, record2},
		}

		mockUsecase := &MockRecordUsecase{
			GetTodayCaloriesFunc: func(ctx context.Context, userID vo.UserID) (*usecase.TodayCaloriesOutput, error) {
				return output, nil
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/records/today", nil)
		c.Set("userID", userIDStr)

		handler.GetToday(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
		}

		var resp dto.TodayCaloriesResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		// 合計カロリーの検証（300+50+800=1150）
		expectedTotal := 1150
		if resp.TotalCalories != expectedTotal {
			t.Errorf("totalCalories = %d, want %d", resp.TotalCalories, expectedTotal)
		}

		// 目標カロリーが設定されていることを確認
		if resp.TargetCalories <= 0 {
			t.Errorf("targetCalories should be positive, got %d", resp.TargetCalories)
		}

		// Recordの数が正しいことを確認
		if len(resp.Records) != 2 {
			t.Errorf("records count = %d, want %d", len(resp.Records), 2)
		}
	})

	t.Run("正常系_記録がない場合", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		// 空のアウトプット
		output := &usecase.TodayCaloriesOutput{
			Date:           time.Now(),
			TotalCalories:  0,
			TargetCalories: 2000,
			Difference:     2000,
			Records:        []*entity.Record{},
		}

		mockUsecase := &MockRecordUsecase{
			GetTodayCaloriesFunc: func(ctx context.Context, userID vo.UserID) (*usecase.TodayCaloriesOutput, error) {
				return output, nil
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/records/today", nil)
		c.Set("userID", userIDStr)

		handler.GetToday(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp dto.TodayCaloriesResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.TotalCalories != 0 {
			t.Errorf("totalCalories = %d, want %d", resp.TotalCalories, 0)
		}

		if len(resp.Records) != 0 {
			t.Errorf("records count = %d, want %d", len(resp.Records), 0)
		}
	})

	t.Run("異常系_認証なし", func(t *testing.T) {
		mockUsecase := &MockRecordUsecase{}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/records/today", nil)
		// userIDを設定しない

		handler.GetToday(c)

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

		mockUsecase := &MockRecordUsecase{
			GetTodayCaloriesFunc: func(ctx context.Context, userID vo.UserID) (*usecase.TodayCaloriesOutput, error) {
				return nil, domainErrors.ErrUserNotFound
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/records/today", nil)
		c.Set("userID", userIDStr)

		handler.GetToday(c)

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

	t.Run("異常系_ユーザー取得でDBエラー", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		mockUsecase := &MockRecordUsecase{
			GetTodayCaloriesFunc: func(ctx context.Context, userID vo.UserID) (*usecase.TodayCaloriesOutput, error) {
				return nil, errors.New("database connection error")
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/records/today", nil)
		c.Set("userID", userIDStr)

		handler.GetToday(c)

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

	t.Run("異常系_Record取得でDBエラー", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		mockUsecase := &MockRecordUsecase{
			GetTodayCaloriesFunc: func(ctx context.Context, userID vo.UserID) (*usecase.TodayCaloriesOutput, error) {
				return nil, errors.New("database connection error")
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/records/today", nil)
		c.Set("userID", userIDStr)

		handler.GetToday(c)

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

func TestRecordHandler_GetStatistics(t *testing.T) {
	t.Run("正常系_週間統計データが取得できる", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		now := time.Now()
		period, _ := vo.NewStatisticsPeriod("week")
		output := &usecase.StatisticsOutput{
			Period:          period,
			TotalDays:       3,
			AverageCalories: vo.ReconstructCalories(2000),
			TargetCalories:  vo.ReconstructCalories(2000),
			AchievedDays:    3,
			OverDays:        0,
			DailyStatistics: []usecase.DailyStatistics{
				{
					Date:           vo.ReconstructEatenAt(now.AddDate(0, 0, -6)),
					TotalCalories:  vo.ReconstructCalories(1800),
					TargetCalories: vo.ReconstructCalories(2000),
					IsAchieved:     true,
					IsOver:         false,
				},
				{
					Date:           vo.ReconstructEatenAt(now.AddDate(0, 0, -5)),
					TotalCalories:  vo.ReconstructCalories(2200),
					TargetCalories: vo.ReconstructCalories(2000),
					IsAchieved:     false,
					IsOver:         true,
				},
				{
					Date:           vo.ReconstructEatenAt(now.AddDate(0, 0, -4)),
					TotalCalories:  vo.ReconstructCalories(2000),
					TargetCalories: vo.ReconstructCalories(2000),
					IsAchieved:     true,
					IsOver:         false,
				},
			},
		}

		mockUsecase := &MockRecordUsecase{
			GetStatisticsFunc: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) (*usecase.StatisticsOutput, error) {
				return output, nil
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/statistics?period=week", nil)
		c.Set("userID", userIDStr)

		handler.GetStatistics(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d, body = %s", w.Code, http.StatusOK, w.Body.String())
		}

		var resp dto.StatisticsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Period != "week" {
			t.Errorf("period = %s, want %s", resp.Period, "week")
		}
		if resp.TotalDays != 3 {
			t.Errorf("totalDays = %d, want %d", resp.TotalDays, 3)
		}
		if len(resp.DailyStatistics) != 3 {
			t.Errorf("dailyStatistics count = %d, want %d", len(resp.DailyStatistics), 3)
		}
		// 平均カロリー: (1800+2200+2000)/3 = 2000
		if resp.AverageCalories != 2000 {
			t.Errorf("averageCalories = %d, want %d", resp.AverageCalories, 2000)
		}
	})

	t.Run("正常系_月間統計データが取得できる", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		period, _ := vo.NewStatisticsPeriod("month")
		output := &usecase.StatisticsOutput{
			Period:          period,
			TotalDays:       0,
			AverageCalories: vo.ReconstructCalories(0),
			TargetCalories:  vo.ReconstructCalories(2000),
			AchievedDays:    0,
			OverDays:        0,
			DailyStatistics: []usecase.DailyStatistics{},
		}

		mockUsecase := &MockRecordUsecase{
			GetStatisticsFunc: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) (*usecase.StatisticsOutput, error) {
				return output, nil
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/statistics?period=month", nil)
		c.Set("userID", userIDStr)

		handler.GetStatistics(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp dto.StatisticsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Period != "month" {
			t.Errorf("period = %s, want %s", resp.Period, "month")
		}
	})

	t.Run("正常系_期間未指定でデフォルトweek", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		period, _ := vo.NewStatisticsPeriod("week")
		output := &usecase.StatisticsOutput{
			Period:          period,
			TotalDays:       0,
			AverageCalories: vo.ReconstructCalories(0),
			TargetCalories:  vo.ReconstructCalories(2000),
			AchievedDays:    0,
			OverDays:        0,
			DailyStatistics: []usecase.DailyStatistics{},
		}

		mockUsecase := &MockRecordUsecase{
			GetStatisticsFunc: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) (*usecase.StatisticsOutput, error) {
				return output, nil
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/statistics", nil)
		c.Set("userID", userIDStr)

		handler.GetStatistics(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp dto.StatisticsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		// 期間未指定の場合はデフォルトでweek
		if resp.Period != "week" {
			t.Errorf("period = %s, want %s", resp.Period, "week")
		}
	})

	t.Run("正常系_データがない場合", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		period, _ := vo.NewStatisticsPeriod("week")
		output := &usecase.StatisticsOutput{
			Period:          period,
			TotalDays:       0,
			AverageCalories: vo.ReconstructCalories(0),
			TargetCalories:  vo.ReconstructCalories(2000),
			AchievedDays:    0,
			OverDays:        0,
			DailyStatistics: []usecase.DailyStatistics{},
		}

		mockUsecase := &MockRecordUsecase{
			GetStatisticsFunc: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) (*usecase.StatisticsOutput, error) {
				return output, nil
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/statistics?period=week", nil)
		c.Set("userID", userIDStr)

		handler.GetStatistics(c)

		if w.Code != http.StatusOK {
			t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
		}

		var resp dto.StatisticsResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.TotalDays != 0 {
			t.Errorf("totalDays = %d, want %d", resp.TotalDays, 0)
		}
		if resp.AverageCalories != 0 {
			t.Errorf("averageCalories = %d, want %d", resp.AverageCalories, 0)
		}
		if len(resp.DailyStatistics) != 0 {
			t.Errorf("dailyStatistics count = %d, want %d", len(resp.DailyStatistics), 0)
		}
	})

	t.Run("異常系_認証なし", func(t *testing.T) {
		mockUsecase := &MockRecordUsecase{}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/statistics?period=week", nil)
		// userIDを設定しない

		handler.GetStatistics(c)

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

	t.Run("異常系_無効な期間パラメータ", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		mockUsecase := &MockRecordUsecase{}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/statistics?period=invalid", nil)
		c.Set("userID", userIDStr)

		handler.GetStatistics(c)

		if w.Code != http.StatusBadRequest {
			t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
		}

		var resp common.ErrorResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if resp.Code != common.CodeValidationError {
			t.Errorf("code = %s, want %s", resp.Code, common.CodeValidationError)
		}
	})

	t.Run("異常系_ユーザーが見つからない", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		mockUsecase := &MockRecordUsecase{
			GetStatisticsFunc: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) (*usecase.StatisticsOutput, error) {
				return nil, domainErrors.ErrUserNotFound
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/statistics?period=week", nil)
		c.Set("userID", userIDStr)

		handler.GetStatistics(c)

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

	t.Run("異常系_ユーザー取得でDBエラー", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		mockUsecase := &MockRecordUsecase{
			GetStatisticsFunc: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) (*usecase.StatisticsOutput, error) {
				return nil, errors.New("database connection error")
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/statistics?period=week", nil)
		c.Set("userID", userIDStr)

		handler.GetStatistics(c)

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

	t.Run("異常系_DailyCalories取得でDBエラー", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"

		mockUsecase := &MockRecordUsecase{
			GetStatisticsFunc: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) (*usecase.StatisticsOutput, error) {
				return nil, errors.New("database connection error")
			},
		}
		handler := record.NewRecordHandler(mockUsecase)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest(http.MethodGet, "/api/v1/statistics?period=week", nil)
		c.Set("userID", userIDStr)

		handler.GetStatistics(c)

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

// domainErrorsのダミー参照（importエラー回避）
var _ = domainErrors.ErrUserNotFound
