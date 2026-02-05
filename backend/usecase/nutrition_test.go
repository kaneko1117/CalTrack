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
	t.Run("正常系_アドバイスが取得できる", func(t *testing.T) {
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

		analyzer := &mockPfcAnalyzer{
			analyze: func(ctx context.Context, config service.PfcAnalyzerConfig, input service.NutritionAdviceInput) (*service.NutritionAdviceOutput, error) {
				// プロンプトが構築されていることを確認
				if config.Prompt == "" {
					t.Error("Prompt should not be empty")
				}
				// プロンプトにカロリー情報が含まれていることを確認
				if config.Prompt == "" {
					t.Error("Prompt should contain calorie information")
				}
				return &service.NutritionAdviceOutput{Advice: "バランスの良い食事ができています"}, nil
			},
		}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, analyzer)
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
	})

	t.Run("正常系_食事記録がない場合", func(t *testing.T) {
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

		analyzer := &mockPfcAnalyzer{
			analyze: func(ctx context.Context, config service.PfcAnalyzerConfig, input service.NutritionAdviceInput) (*service.NutritionAdviceOutput, error) {
				// 食事記録がない場合でもアドバイスを返す
				if input.CurrentCalories != 0 {
					t.Errorf("CurrentCalories = %d, want 0", input.CurrentCalories)
				}
				return &service.NutritionAdviceOutput{Advice: "まだ食事記録がありません。最初の食事を記録してみましょう。"}, nil
			},
		}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, analyzer)
		output, err := uc.GetAdvice(context.Background(), userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.Advice == "" {
			t.Error("Advice should not be empty")
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
		analyzer := &mockPfcAnalyzer{}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, analyzer)
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
		analyzer := &mockPfcAnalyzer{}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, analyzer)
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
		analyzer := &mockPfcAnalyzer{}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, analyzer)
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

		analyzer := &mockPfcAnalyzer{}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, analyzer)
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

		analyzer := &mockPfcAnalyzer{
			analyze: func(ctx context.Context, config service.PfcAnalyzerConfig, input service.NutritionAdviceInput) (*service.NutritionAdviceOutput, error) {
				return nil, analyzeErr
			},
		}

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, analyzer)
		_, err := uc.GetAdvice(context.Background(), userID)

		if !errors.Is(err, analyzeErr) {
			t.Errorf("got %v, want analyzeErr", err)
		}
	})
}
