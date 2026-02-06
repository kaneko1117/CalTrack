package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/repository"
	"caltrack/domain/vo"
	"caltrack/usecase"
	"caltrack/usecase/service"
)

// mockAdviceCacheRepository はAdviceCacheRepositoryのモック実装
type mockAdviceCacheRepository struct {
	save                 func(ctx context.Context, cache *entity.AdviceCache) error
	findByUserIDAndDate  func(ctx context.Context, userID vo.UserID, date time.Time) (*entity.AdviceCache, error)
	deleteByUserIDAndDate func(ctx context.Context, userID vo.UserID, date time.Time) error
}

func (m *mockAdviceCacheRepository) Save(ctx context.Context, cache *entity.AdviceCache) error {
	if m.save != nil {
		return m.save(ctx, cache)
	}
	return nil
}

func (m *mockAdviceCacheRepository) FindByUserIDAndDate(ctx context.Context, userID vo.UserID, date time.Time) (*entity.AdviceCache, error) {
	if m.findByUserIDAndDate != nil {
		return m.findByUserIDAndDate(ctx, userID, date)
	}
	return nil, nil
}

func (m *mockAdviceCacheRepository) DeleteByUserIDAndDate(ctx context.Context, userID vo.UserID, date time.Time) error {
	if m.deleteByUserIDAndDate != nil {
		return m.deleteByUserIDAndDate(ctx, userID, date)
	}
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
	return &service.NutritionAdviceOutput{Advice: "テストアドバイス"}, nil
}

// mockNutritionUserRepository はUserRepositoryのモック実装（Nutrition用）
type mockNutritionUserRepository struct {
	findByID func(ctx context.Context, id vo.UserID) (*entity.User, error)
}

func (m *mockNutritionUserRepository) Save(ctx context.Context, user *entity.User) error {
	return nil
}

func (m *mockNutritionUserRepository) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	return nil, nil
}

func (m *mockNutritionUserRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	return false, nil
}

func (m *mockNutritionUserRepository) FindByID(ctx context.Context, id vo.UserID) (*entity.User, error) {
	if m.findByID != nil {
		return m.findByID(ctx, id)
	}
	return nil, nil
}

func (m *mockNutritionUserRepository) Update(ctx context.Context, user *entity.User) error {
	return nil
}

// mockNutritionRecordRepository はRecordRepositoryのモック実装（Nutrition用）
type mockNutritionRecordRepository struct {
	findByUserIDAndDateRange func(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error)
}

func (m *mockNutritionRecordRepository) Save(ctx context.Context, record *entity.Record) error {
	return nil
}

func (m *mockNutritionRecordRepository) FindByUserIDAndDateRange(ctx context.Context, userID vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
	if m.findByUserIDAndDateRange != nil {
		return m.findByUserIDAndDateRange(ctx, userID, startTime, endTime)
	}
	return nil, nil
}

func (m *mockNutritionRecordRepository) GetDailyCalories(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
	return nil, nil
}

// mockNutritionRecordPfcRepository はRecordPfcRepositoryのモック実装（Nutrition用）
type mockNutritionRecordPfcRepository struct {
	findByRecordIDs func(ctx context.Context, recordIDs []vo.RecordID) ([]*entity.RecordPfc, error)
}

func (m *mockNutritionRecordPfcRepository) Save(ctx context.Context, recordPfc *entity.RecordPfc) error {
	return nil
}

func (m *mockNutritionRecordPfcRepository) FindByRecordID(ctx context.Context, recordID vo.RecordID) (*entity.RecordPfc, error) {
	return nil, nil
}

func (m *mockNutritionRecordPfcRepository) FindByRecordIDs(ctx context.Context, recordIDs []vo.RecordID) ([]*entity.RecordPfc, error) {
	if m.findByRecordIDs != nil {
		return m.findByRecordIDs(ctx, recordIDs)
	}
	return nil, nil
}

// validUserForNutrition はNutrition用テストのための有効なUserを生成する
func validUserForNutrition(t *testing.T, userID vo.UserID) *entity.User {
	t.Helper()
	user, err := entity.ReconstructUser(
		userID.String(),
		"test@example.com",
		"hashedpassword",
		"testuser",
		70.0,
		175.0,
		time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
		"male",
		"moderate",
		time.Now(),
		time.Now(),
	)
	if err != nil {
		t.Fatalf("failed to create valid user: %v", err)
	}
	return user
}

func TestNutritionUsecase_GetAdvice(t *testing.T) {
	t.Run("正常系_今日の記録がない場合は固定文言が返される", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForNutrition(t, userID)

		userRepo := &mockNutritionUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}

		recordRepo := &mockNutritionRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return []*entity.Record{}, nil
			},
		}

		recordPfcRepo := &mockNutritionRecordPfcRepository{}
		adviceCacheRepo := &mockAdviceCacheRepository{}
		analyzer := &mockPfcAnalyzer{}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		output, err := uc.GetAdvice(context.Background(), userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.Advice != usecase.NoRecordAdviceMessage {
			t.Errorf("Advice = %s, want %s", output.Advice, usecase.NoRecordAdviceMessage)
		}
	})

	t.Run("正常系_キャッシュがある場合はキャッシュが返される", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForNutrition(t, userID)

		// 今日のRecordを作成
		record1, _ := entity.NewRecord(userID, time.Now())
		_ = record1.AddItem("朝食", 300)
		records := []*entity.Record{record1}

		cachedAdvice := "キャッシュされたアドバイス"
		cache := entity.NewAdviceCache(userID, time.Now(), cachedAdvice)

		userRepo := &mockNutritionUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}

		recordRepo := &mockNutritionRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return records, nil
			},
		}

		recordPfcRepo := &mockNutritionRecordPfcRepository{}

		adviceCacheRepo := &mockAdviceCacheRepository{
			findByUserIDAndDate: func(ctx context.Context, uid vo.UserID, date time.Time) (*entity.AdviceCache, error) {
				return cache, nil
			},
		}

		// AI呼び出しは行われないはず
		analyzerCalled := false
		analyzer := &mockPfcAnalyzer{
			analyze: func(ctx context.Context, config service.PfcAnalyzerConfig, input service.NutritionAdviceInput) (*service.NutritionAdviceOutput, error) {
				analyzerCalled = true
				return &service.NutritionAdviceOutput{Advice: "新しいアドバイス"}, nil
			},
		}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		output, err := uc.GetAdvice(context.Background(), userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.Advice != cachedAdvice {
			t.Errorf("Advice = %s, want %s", output.Advice, cachedAdvice)
		}

		if analyzerCalled {
			t.Error("Analyzer should not be called when cache exists")
		}
	})

	t.Run("正常系_キャッシュがない場合はAI呼び出し後にキャッシュ保存される", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForNutrition(t, userID)

		// 今日のRecordを作成
		record1, _ := entity.NewRecord(userID, time.Now())
		_ = record1.AddItem("朝食", 300)

		record2, _ := entity.NewRecord(userID, time.Now())
		_ = record2.AddItem("昼食", 500)

		records := []*entity.Record{record1, record2}

		// RecordPfcを作成
		recordPfc1 := entity.NewRecordPfc(record1.ID(), 15.0, 10.0, 40.0)
		recordPfc2 := entity.NewRecordPfc(record2.ID(), 20.0, 15.0, 60.0)

		userRepo := &mockNutritionUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}

		recordRepo := &mockNutritionRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return records, nil
			},
		}

		recordPfcRepo := &mockNutritionRecordPfcRepository{
			findByRecordIDs: func(ctx context.Context, recordIDs []vo.RecordID) ([]*entity.RecordPfc, error) {
				return []*entity.RecordPfc{recordPfc1, recordPfc2}, nil
			},
		}

		var savedCache *entity.AdviceCache
		adviceCacheRepo := &mockAdviceCacheRepository{
			findByUserIDAndDate: func(ctx context.Context, uid vo.UserID, date time.Time) (*entity.AdviceCache, error) {
				return nil, nil // キャッシュなし
			},
			save: func(ctx context.Context, cache *entity.AdviceCache) error {
				savedCache = cache
				return nil
			},
		}

		analyzerCalled := false
		analyzer := &mockPfcAnalyzer{
			analyze: func(ctx context.Context, config service.PfcAnalyzerConfig, input service.NutritionAdviceInput) (*service.NutritionAdviceOutput, error) {
				analyzerCalled = true
				// プロンプトが構築されていることを確認
				if config.Prompt == "" {
					t.Error("Prompt should not be empty")
				}
				return &service.NutritionAdviceOutput{Advice: "バランスの良い食事ができています"}, nil
			},
		}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		output, err := uc.GetAdvice(context.Background(), userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.Advice == "" {
			t.Error("Advice should not be empty")
		}

		if output.Advice != "バランスの良い食事ができています" {
			t.Errorf("Advice = %s, want %s", output.Advice, "バランスの良い食事ができています")
		}

		if !analyzerCalled {
			t.Error("Analyzer should be called when cache does not exist")
		}

		if savedCache == nil {
			t.Error("Cache should be saved")
		}

		if savedCache != nil && savedCache.Advice() != "バランスの良い食事ができています" {
			t.Errorf("Saved cache advice = %s, want %s", savedCache.Advice(), "バランスの良い食事ができています")
		}
	})


	t.Run("異常系_ユーザーが存在しない", func(t *testing.T) {
		userID := vo.NewUserID()

		userRepo := &mockNutritionUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, nil
			},
		}

		recordRepo := &mockNutritionRecordRepository{}
		recordPfcRepo := &mockNutritionRecordPfcRepository{}
		adviceCacheRepo := &mockAdviceCacheRepository{}
		analyzer := &mockPfcAnalyzer{}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		_, err := uc.GetAdvice(context.Background(), userID)

		if !errors.Is(err, domainErrors.ErrUserNotFound) {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_ユーザー取得時にエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		repoErr := errors.New("db error")

		userRepo := &mockNutritionUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, repoErr
			},
		}

		recordRepo := &mockNutritionRecordRepository{}
		recordPfcRepo := &mockNutritionRecordPfcRepository{}
		adviceCacheRepo := &mockAdviceCacheRepository{}
		analyzer := &mockPfcAnalyzer{}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		_, err := uc.GetAdvice(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_Record取得時にエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForNutrition(t, userID)
		repoErr := errors.New("db error")

		userRepo := &mockNutritionUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}

		recordRepo := &mockNutritionRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return nil, repoErr
			},
		}

		recordPfcRepo := &mockNutritionRecordPfcRepository{}
		adviceCacheRepo := &mockAdviceCacheRepository{}
		analyzer := &mockPfcAnalyzer{}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		_, err := uc.GetAdvice(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_RecordPfc取得時にエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForNutrition(t, userID)
		repoErr := errors.New("db error")

		record1, _ := entity.NewRecord(userID, time.Now())
		_ = record1.AddItem("朝食", 300)

		userRepo := &mockNutritionUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}

		recordRepo := &mockNutritionRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return []*entity.Record{record1}, nil
			},
		}

		recordPfcRepo := &mockNutritionRecordPfcRepository{
			findByRecordIDs: func(ctx context.Context, recordIDs []vo.RecordID) ([]*entity.RecordPfc, error) {
				return nil, repoErr
			},
		}

		adviceCacheRepo := &mockAdviceCacheRepository{
			findByUserIDAndDate: func(ctx context.Context, uid vo.UserID, date time.Time) (*entity.AdviceCache, error) {
				return nil, nil
			},
		}

		analyzer := &mockPfcAnalyzer{}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		_, err := uc.GetAdvice(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_PfcAnalyzer実行時にエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForNutrition(t, userID)
		analyzeErr := errors.New("analyzer error")

		record1, _ := entity.NewRecord(userID, time.Now())
		_ = record1.AddItem("朝食", 300)

		userRepo := &mockNutritionUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}

		recordRepo := &mockNutritionRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return []*entity.Record{record1}, nil
			},
		}

		recordPfcRepo := &mockNutritionRecordPfcRepository{
			findByRecordIDs: func(ctx context.Context, recordIDs []vo.RecordID) ([]*entity.RecordPfc, error) {
				return []*entity.RecordPfc{}, nil
			},
		}

		adviceCacheRepo := &mockAdviceCacheRepository{
			findByUserIDAndDate: func(ctx context.Context, uid vo.UserID, date time.Time) (*entity.AdviceCache, error) {
				return nil, nil
			},
		}

		analyzer := &mockPfcAnalyzer{
			analyze: func(ctx context.Context, config service.PfcAnalyzerConfig, input service.NutritionAdviceInput) (*service.NutritionAdviceOutput, error) {
				return nil, analyzeErr
			},
		}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		_, err := uc.GetAdvice(context.Background(), userID)

		if !errors.Is(err, analyzeErr) {
			t.Errorf("got %v, want analyzeErr", err)
		}
	})
}
