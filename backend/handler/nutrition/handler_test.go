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

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/repository"
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

// mockUserRepository はUserRepositoryのモック実装
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

// mockRecordRepository はRecordRepositoryのモック実装
type mockRecordRepository struct {
	findByUserIDAndDateRange func(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error)
}

func (m *mockRecordRepository) Save(ctx context.Context, record *entity.Record) error {
	return nil
}

func (m *mockRecordRepository) FindByUserIDAndDateRange(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
	if m.findByUserIDAndDateRange != nil {
		return m.findByUserIDAndDateRange(ctx, userID, startTime, endTime)
	}
	return nil, nil
}

func (m *mockRecordRepository) GetDailyCalories(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
	return nil, nil
}

// mockRecordPfcRepository はRecordPfcRepositoryのモック実装
type mockRecordPfcRepository struct{}

func (m *mockRecordPfcRepository) Save(ctx context.Context, recordPfc *entity.RecordPfc) error {
	return nil
}

func (m *mockRecordPfcRepository) FindByRecordID(ctx context.Context, recordID vo.RecordID) (*entity.RecordPfc, error) {
	return nil, nil
}

func (m *mockRecordPfcRepository) FindByRecordIDs(ctx context.Context, recordIDs []vo.RecordID) ([]*entity.RecordPfc, error) {
	return nil, nil
}

// mockAdviceCacheRepository はAdviceCacheRepositoryのモック実装
type mockAdviceCacheRepository struct{}

func (m *mockAdviceCacheRepository) Save(ctx context.Context, cache *entity.AdviceCache) error {
	return nil
}

func (m *mockAdviceCacheRepository) FindByUserIDAndDate(ctx context.Context, userID vo.UserID, date time.Time) (*entity.AdviceCache, error) {
	return nil, nil
}

func (m *mockAdviceCacheRepository) DeleteByUserIDAndDate(ctx context.Context, userID vo.UserID, date time.Time) error {
	return nil
}

// mockPfcAnalyzer はPfcAnalyzerのモック実装
type mockPfcAnalyzer struct {
	analyze func(ctx context.Context, config service.PfcAnalyzerConfig, input service.NutritionAdviceInput) (*service.NutritionAdviceOutput, error)
}

func (m *mockPfcAnalyzer) Analyze(ctx context.Context, config service.PfcAnalyzerConfig, input service.NutritionAdviceInput) (*service.NutritionAdviceOutput, error) {
	if m.analyze != nil {
		return m.analyze(ctx, config, input)
	}
	return &service.NutritionAdviceOutput{Advice: "モックアドバイス"}, nil
}

// setupHandler はテスト用のハンドラをセットアップする
func setupHandler(userRepo repository.UserRepository, recordRepo repository.RecordRepository, analyzer service.PfcAnalyzer) *nutrition.NutritionHandler {
	recordPfcRepo := &mockRecordPfcRepository{}
	adviceCacheRepo := &mockAdviceCacheRepository{}
	uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
	return nutrition.NewNutritionHandler(uc)
}

func TestNutritionHandler_GetAdvice(t *testing.T) {
	t.Run("正常系_アドバイス取得成功", func(t *testing.T) {
		userIDStr := "550e8400-e29b-41d4-a716-446655440000"
		userID := vo.ReconstructUserID(userIDStr)
		birthDate := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
		now := time.Now()

		// ダミーのRecordを作成
		record, _ := entity.NewRecord(userID, now)
		_ = record.AddItem("テスト食事", 500)

		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return entity.ReconstructUser(
					userIDStr,
					"test@example.com",
					"$2a$10$dummy_hash",
					"テストユーザー",
					70.0,
					170.0,
					birthDate,
					"male",
					"moderate",
					now,
					now,
				)
			},
		}

		recordRepo := &mockRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return []*entity.Record{record}, nil
			},
		}

		analyzer := &mockPfcAnalyzer{
			analyze: func(ctx context.Context, config service.PfcAnalyzerConfig, input service.NutritionAdviceInput) (*service.NutritionAdviceOutput, error) {
				return &service.NutritionAdviceOutput{
					Advice: "バランスの良い食事を心がけましょう。",
				}, nil
			},
		}

		handler := setupHandler(userRepo, recordRepo, analyzer)

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
		userRepo := &mockUserRepository{}
		analyzer := &mockPfcAnalyzer{}
		recordRepo := &mockRecordRepository{}
		handler := setupHandler(userRepo, recordRepo, analyzer)

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

		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, domainErrors.ErrUserNotFound
			},
		}
		analyzer := &mockPfcAnalyzer{}
		recordRepo := &mockRecordRepository{}
		handler := setupHandler(userRepo, recordRepo, analyzer)

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

		userRepo := &mockUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, errors.New("database connection error")
			},
		}
		analyzer := &mockPfcAnalyzer{}
		recordRepo := &mockRecordRepository{}
		handler := setupHandler(userRepo, recordRepo, analyzer)

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
