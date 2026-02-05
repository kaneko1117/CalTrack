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
	getDailyCalories         func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) ([]repository.DailyCalories, error)
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

func (m *mockRecordRepository) GetDailyCalories(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
	if m.getDailyCalories != nil {
		return m.getDailyCalories(ctx, userID, period)
	}
	return nil, nil
}

// mockTransactionManager はTransactionManagerのモック実装
type mockTransactionManager struct{}

func (m *mockTransactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// mockRecordPfcRepository はRecordPfcRepositoryのモック実装
type mockRecordPfcRepository struct {
	save            func(ctx context.Context, recordPfc *entity.RecordPfc) error
	findByRecordID  func(ctx context.Context, recordID vo.RecordID) (*entity.RecordPfc, error)
	findByRecordIDs func(ctx context.Context, recordIDs []vo.RecordID) ([]*entity.RecordPfc, error)
}

func (m *mockRecordPfcRepository) Save(ctx context.Context, recordPfc *entity.RecordPfc) error {
	if m.save != nil {
		return m.save(ctx, recordPfc)
	}
	return nil
}

func (m *mockRecordPfcRepository) FindByRecordID(ctx context.Context, recordID vo.RecordID) (*entity.RecordPfc, error) {
	if m.findByRecordID != nil {
		return m.findByRecordID(ctx, recordID)
	}
	return nil, nil
}

func (m *mockRecordPfcRepository) FindByRecordIDs(ctx context.Context, recordIDs []vo.RecordID) ([]*entity.RecordPfc, error) {
	if m.findByRecordIDs != nil {
		return m.findByRecordIDs(ctx, recordIDs)
	}
	return nil, nil
}

// mockUserRepository はUserRepositoryのモック実装（Handler用）
type mockUserRepository struct {
	findByID func(ctx context.Context, id vo.UserID) (*entity.User, error)
}

func (m *mockUserRepository) Save(ctx context.Context, user *entity.User) error {
	return nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	return nil, nil
}

func (m *mockUserRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	return false, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id vo.UserID) (*entity.User, error) {
	if m.findByID != nil {
		return m.findByID(ctx, id)
	}
	return nil, nil
}

// setupHandler はテスト用のハンドラをセットアップする
func setupHandler(recordRepo repository.RecordRepository) *record.RecordHandler {
	recordPfcRepo := &mockRecordPfcRepository{}
	userRepo := &mockUserRepository{}
	txManager := &mockTransactionManager{}
	uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, txManager)
	return record.NewRecordHandler(uc)
}

// setupHandlerWithUserRepo はUserRepositoryを指定してテスト用のハンドラをセットアップする
func setupHandlerWithUserRepo(recordRepo repository.RecordRepository, userRepo repository.UserRepository) *record.RecordHandler {
	recordPfcRepo := &mockRecordPfcRepository{}
	txManager := &mockTransactionManager{}
	uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, txManager)
	return record.NewRecordHandler(uc)
}

// createTestUser はテスト用のユーザーを作成する
func createTestUser(userIDStr string) *entity.User {
	birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	user, _ := entity.ReconstructUser(
		userIDStr,
		"test@example.com",
		"$2a$10$dummy_hash_value_for_testing",
		"テストユーザー",
		70.0,  // weight
		170.0, // height
		birthDate,
		"male",
		"moderate",
		time.Now(),
		time.Now(),
	)
	return user
}

// createTestRecord はテスト用のRecordを作成する
func createTestRecord(userIDStr string, eatenAt time.Time, items []struct {
	name     string
	calories int
}) *entity.Record {
	recordItems := make([]entity.RecordItem, len(items))
	recordIDStr := "record-" + eatenAt.Format("20060102150405")
	for i, item := range items {
		recordItems[i] = *entity.ReconstructRecordItem(
			"item-"+recordIDStr+"-"+item.name,
			recordIDStr,
			item.name,
			item.calories,
		)
	}
	return entity.ReconstructRecord(
		recordIDStr,
		userIDStr,
		eatenAt,
		time.Now(),
		recordItems,
	)
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

func TestRecordHandler_GetToday(t *testing.T) {
	t.Run("正常系_今日のカロリー情報が取得できる", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"
		testUser := createTestUser(userIDStr)

		// 今日のRecordを作成
		now := time.Now()
		testRecords := []*entity.Record{
			createTestRecord(userIDStr, now.Add(-2*time.Hour), []struct {
				name     string
				calories int
			}{
				{"朝食：パン", 300},
				{"朝食：コーヒー", 50},
			}),
			createTestRecord(userIDStr, now.Add(-1*time.Hour), []struct {
				name     string
				calories int
			}{
				{"昼食：ラーメン", 800},
			}),
		}

		recordRepo := &mockRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return testRecords, nil
			},
		}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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
		testUser := createTestUser(userIDStr)

		recordRepo := &mockRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return []*entity.Record{}, nil
			},
		}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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
		recordRepo := &mockRecordRepository{}
		userRepo := &mockUserRepository{}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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

		recordRepo := &mockRecordRepository{}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, nil // ユーザーが見つからない
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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

		recordRepo := &mockRecordRepository{}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, errors.New("database connection error")
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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
		testUser := createTestUser(userIDStr)

		recordRepo := &mockRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return nil, errors.New("database connection error")
			},
		}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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
		testUser := createTestUser(userIDStr)

		// 日別カロリーデータを作成
		now := time.Now()
		dailyCalories := []repository.DailyCalories{
			{Date: vo.ReconstructEatenAt(now.AddDate(0, 0, -6)), Calories: vo.ReconstructCalories(1800)},
			{Date: vo.ReconstructEatenAt(now.AddDate(0, 0, -5)), Calories: vo.ReconstructCalories(2200)},
			{Date: vo.ReconstructEatenAt(now.AddDate(0, 0, -4)), Calories: vo.ReconstructCalories(2000)},
		}

		recordRepo := &mockRecordRepository{
			getDailyCalories: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
				return dailyCalories, nil
			},
		}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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
		testUser := createTestUser(userIDStr)

		recordRepo := &mockRecordRepository{
			getDailyCalories: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
				return []repository.DailyCalories{}, nil
			},
		}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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
		testUser := createTestUser(userIDStr)

		recordRepo := &mockRecordRepository{
			getDailyCalories: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
				return []repository.DailyCalories{}, nil
			},
		}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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
		testUser := createTestUser(userIDStr)

		recordRepo := &mockRecordRepository{
			getDailyCalories: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
				return []repository.DailyCalories{}, nil
			},
		}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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
		recordRepo := &mockRecordRepository{}
		userRepo := &mockUserRepository{}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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

		recordRepo := &mockRecordRepository{}
		userRepo := &mockUserRepository{}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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

		recordRepo := &mockRecordRepository{}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, nil // ユーザーが見つからない
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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

		recordRepo := &mockRecordRepository{}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, errors.New("database connection error")
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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
		testUser := createTestUser(userIDStr)

		recordRepo := &mockRecordRepository{
			getDailyCalories: func(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
				return nil, errors.New("database connection error")
			},
		}
		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return testUser, nil
			},
		}
		handler := setupHandlerWithUserRepo(recordRepo, userRepo)

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
