package record_test

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
	"caltrack/domain/repository"
	"caltrack/domain/vo"
	"caltrack/handler/common"
	"caltrack/handler/record"
	"caltrack/handler/record/dto"
	"caltrack/usecase"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// mockRecordRepository はRecordRepositoryのモック実装
type mockRecordRepository struct {
	save                     func(ctx context.Context, record *entity.Record) error
	findByUserIDAndDateRange func(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error)
}

func (m *mockRecordRepository) Save(ctx context.Context, record *entity.Record) error {
	return m.save(ctx, record)
}

func (m *mockRecordRepository) FindByUserIDAndDateRange(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
	if m.findByUserIDAndDateRange != nil {
		return m.findByUserIDAndDateRange(ctx, userID, startTime, endTime)
	}
	return nil, nil
}

// mockTransactionManager はTransactionManagerのモック実装
type mockTransactionManager struct{}

func (m *mockTransactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// setupHandler はテスト用のハンドラをセットアップする
func setupHandler(recordRepo repository.RecordRepository) *record.RecordHandler {
	txManager := &mockTransactionManager{}
	uc := usecase.NewRecordUsecase(recordRepo, txManager)
	return record.NewRecordHandler(uc)
}

func TestRecordHandler_Create(t *testing.T) {
	t.Run("正常系_記録が作成される", func(t *testing.T) {
		recordRepo := &mockRecordRepository{
			save: func(ctx context.Context, record *entity.Record) error {
				return nil
			},
		}
		handler := setupHandler(recordRepo)

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
		recordRepo := &mockRecordRepository{}
		handler := setupHandler(recordRepo)

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
		recordRepo := &mockRecordRepository{}
		handler := setupHandler(recordRepo)

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
		recordRepo := &mockRecordRepository{}
		handler := setupHandler(recordRepo)

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
		recordRepo := &mockRecordRepository{}
		handler := setupHandler(recordRepo)

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
		recordRepo := &mockRecordRepository{}
		handler := setupHandler(recordRepo)

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
		recordRepo := &mockRecordRepository{}
		handler := setupHandler(recordRepo)

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
		recordRepo := &mockRecordRepository{}
		handler := setupHandler(recordRepo)

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
		recordRepo := &mockRecordRepository{
			save: func(ctx context.Context, record *entity.Record) error {
				return context.DeadlineExceeded
			},
		}
		handler := setupHandler(recordRepo)

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
