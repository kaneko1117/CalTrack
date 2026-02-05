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
)

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

// mockRecordUserRepository はUserRepositoryのモック実装（Record用）
type mockRecordUserRepository struct {
	findByID func(ctx context.Context, id vo.UserID) (*entity.User, error)
}

func (m *mockRecordUserRepository) Save(ctx context.Context, user *entity.User) error {
	return nil
}

func (m *mockRecordUserRepository) FindByEmail(ctx context.Context, email vo.Email) (*entity.User, error) {
	return nil, nil
}

func (m *mockRecordUserRepository) ExistsByEmail(ctx context.Context, email vo.Email) (bool, error) {
	return false, nil
}

func (m *mockRecordUserRepository) FindByID(ctx context.Context, id vo.UserID) (*entity.User, error) {
	if m.findByID != nil {
		return m.findByID(ctx, id)
	}
	return nil, nil
}

// mockRecordAdviceCacheRepository はAdviceCacheRepositoryのモック実装
type mockRecordAdviceCacheRepository struct {
	deleteByUserIDAndDate func(ctx context.Context, userID vo.UserID, date time.Time) error
}

func (m *mockRecordAdviceCacheRepository) Save(ctx context.Context, cache *entity.AdviceCache) error {
	return nil
}

func (m *mockRecordAdviceCacheRepository) FindByUserIDAndDate(ctx context.Context, userID vo.UserID, date time.Time) (*entity.AdviceCache, error) {
	return nil, nil
}

func (m *mockRecordAdviceCacheRepository) DeleteByUserIDAndDate(ctx context.Context, userID vo.UserID, date time.Time) error {
	if m.deleteByUserIDAndDate != nil {
		return m.deleteByUserIDAndDate(ctx, userID, date)
	}
	return nil
}

// mockRecordTransactionManager はTransactionManagerのモック実装
type mockRecordTransactionManager struct{}

func (m *mockRecordTransactionManager) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

// validRecord はテスト用の有効なRecordを生成する
func validRecord(t *testing.T) *entity.Record {
	t.Helper()
	userID := vo.NewUserID()
	record, err := entity.NewRecord(userID, time.Now())
	if err != nil {
		t.Fatalf("failed to create valid record: %v", err)
	}
	return record
}

func TestRecordUsecase_Create(t *testing.T) {
	t.Run("正常系_記録が保存されキャッシュが無効化される", func(t *testing.T) {
		var savedRecord *entity.Record
		var savedRecordPfc *entity.RecordPfc
		cacheDeleted := false
		recordRepo := &mockRecordRepository{
			save: func(ctx context.Context, record *entity.Record) error {
				savedRecord = record
				return nil
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{
			save: func(ctx context.Context, recordPfc *entity.RecordPfc) error {
				savedRecordPfc = recordPfc
				return nil
			},
		}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{
			deleteByUserIDAndDate: func(ctx context.Context, userID vo.UserID, date time.Time) error {
				cacheDeleted = true
				return nil
			},
		}
		userRepo := &mockRecordUserRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		record := validRecord(t)
		recordPfc := entity.NewRecordPfc(record.ID(), 20.0, 10.0, 30.0)
		err := uc.Create(context.Background(), record, recordPfc)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if savedRecord == nil {
			t.Error("record should be saved")
		}
		if savedRecord.ID().String() != record.ID().String() {
			t.Errorf("saved record ID = %s, want %s", savedRecord.ID().String(), record.ID().String())
		}
		if savedRecordPfc == nil {
			t.Error("recordPfc should be saved")
		}
		if !cacheDeleted {
			t.Error("cache should be deleted")
		}
	})

	t.Run("異常系_保存時にエラーが発生", func(t *testing.T) {
		saveErr := errors.New("save error")
		recordRepo := &mockRecordRepository{
			save: func(ctx context.Context, record *entity.Record) error {
				return saveErr
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		userRepo := &mockRecordUserRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		record := validRecord(t)
		recordPfc := entity.NewRecordPfc(record.ID(), 20.0, 10.0, 30.0)
		err := uc.Create(context.Background(), record, recordPfc)

		if !errors.Is(err, saveErr) {
			t.Errorf("got %v, want saveErr", err)
		}
	})
}

// validUserForRecord はRecord用テストのための有効なUserを生成する
func validUserForRecord(t *testing.T, userID vo.UserID) *entity.User {
	t.Helper()
	// ReconstructUserを使用してuserIDを指定できるようにする
	user, err := entity.ReconstructUser(
		userID.String(),
		"test@example.com",
		"hashedpassword",
		"testuser",
		70.5,
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

func TestRecordUsecase_GetTodayCalories(t *testing.T) {
	t.Run("正常系_今日のカロリー情報を取得", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		// 今日のRecordを作成
		record1, _ := entity.NewRecord(userID, time.Now())
		_ = record1.AddItem("朝食", 300)

		record2, _ := entity.NewRecord(userID, time.Now())
		_ = record2.AddItem("昼食", 500)
		_ = record2.AddItem("おやつ", 100)

		records := []*entity.Record{record1, record2}

		recordRepo := &mockRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return records, nil
			},
		}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		output, err := uc.GetTodayCalories(context.Background(), userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 合計カロリーの検証（300 + 500 + 100 = 900）
		if output.TotalCalories != 900 {
			t.Errorf("TotalCalories = %d, want 900", output.TotalCalories)
		}

		// 目標カロリーの検証
		expectedTargetCalories := user.CalculateTargetCalories()
		if output.TargetCalories != expectedTargetCalories {
			t.Errorf("TargetCalories = %d, want %d", output.TargetCalories, expectedTargetCalories)
		}

		// 差分の検証（目標 - 実績）
		expectedDifference := expectedTargetCalories - 900
		if output.Difference != expectedDifference {
			t.Errorf("Difference = %d, want %d", output.Difference, expectedDifference)
		}

		// Record数の検証
		if len(output.Records) != 2 {
			t.Errorf("len(Records) = %d, want 2", len(output.Records))
		}
	})

	t.Run("正常系_記録が0件の場合", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		recordRepo := &mockRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return []*entity.Record{}, nil
			},
		}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		output, err := uc.GetTodayCalories(context.Background(), userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 合計カロリーは0
		if output.TotalCalories != 0 {
			t.Errorf("TotalCalories = %d, want 0", output.TotalCalories)
		}

		// Record数は0
		if len(output.Records) != 0 {
			t.Errorf("len(Records) = %d, want 0", len(output.Records))
		}
	})

	t.Run("異常系_ユーザーが存在しない", func(t *testing.T) {
		userID := vo.NewUserID()

		recordRepo := &mockRecordRepository{}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, nil // ユーザーが存在しない
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		_, err := uc.GetTodayCalories(context.Background(), userID)

		if !errors.Is(err, domainErrors.ErrUserNotFound) {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_ユーザー取得時にエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		repoErr := errors.New("db error")

		recordRepo := &mockRecordRepository{}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, repoErr
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		_, err := uc.GetTodayCalories(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_Record取得時にエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		repoErr := errors.New("db error")

		recordRepo := &mockRecordRepository{
			findByUserIDAndDateRange: func(ctx context.Context, uid vo.UserID, startTime, endTime time.Time) ([]*entity.Record, error) {
				return nil, repoErr
			},
		}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		_, err := uc.GetTodayCalories(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})
}

func TestRecordUsecase_GetStatistics(t *testing.T) {
	t.Run("正常系_週間統計データを取得", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		targetCalories := user.CalculateTargetCalories()

		// 7日分のデータを作成（達成:3日、超過:2日、未達成:2日）
		// 達成条件: 80%以上100%以下
		// 超過条件: 100%超
		// 未達成: 80%未満
		now := time.Now()
		dailyCalories := []repository.DailyCalories{
			// 達成（85%）
			{Date: vo.ReconstructEatenAt(now.AddDate(0, 0, -6)), Calories: vo.ReconstructCalories(targetCalories * 85 / 100)},
			// 達成（90%）
			{Date: vo.ReconstructEatenAt(now.AddDate(0, 0, -5)), Calories: vo.ReconstructCalories(targetCalories * 90 / 100)},
			// 達成（100%ちょうど）
			{Date: vo.ReconstructEatenAt(now.AddDate(0, 0, -4)), Calories: vo.ReconstructCalories(targetCalories)},
			// 超過（120%）
			{Date: vo.ReconstructEatenAt(now.AddDate(0, 0, -3)), Calories: vo.ReconstructCalories(targetCalories * 120 / 100)},
			// 超過（150%）
			{Date: vo.ReconstructEatenAt(now.AddDate(0, 0, -2)), Calories: vo.ReconstructCalories(targetCalories * 150 / 100)},
			// 未達成（50%）
			{Date: vo.ReconstructEatenAt(now.AddDate(0, 0, -1)), Calories: vo.ReconstructCalories(targetCalories * 50 / 100)},
			// 未達成（30%）
			{Date: vo.ReconstructEatenAt(now), Calories: vo.ReconstructCalories(targetCalories * 30 / 100)},
		}

		period, _ := vo.NewStatisticsPeriod("week")

		recordRepo := &mockRecordRepository{
			getDailyCalories: func(ctx context.Context, uid vo.UserID, p vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
				return dailyCalories, nil
			},
		}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		output, err := uc.GetStatistics(context.Background(), userID, period)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 期間の検証
		if output.Period.String() != "week" {
			t.Errorf("Period = %s, want week", output.Period.String())
		}

		// 目標カロリーの検証
		if output.TargetCalories.Value() != targetCalories {
			t.Errorf("TargetCalories = %d, want %d", output.TargetCalories.Value(), targetCalories)
		}

		// 日数の検証
		if output.TotalDays != 7 {
			t.Errorf("TotalDays = %d, want 7", output.TotalDays)
		}

		// 達成日数の検証（80%～100%: 3日）
		if output.AchievedDays != 3 {
			t.Errorf("AchievedDays = %d, want 3", output.AchievedDays)
		}

		// 超過日数の検証（100%超: 2日）
		if output.OverDays != 2 {
			t.Errorf("OverDays = %d, want 2", output.OverDays)
		}

		// DailyStatistics数の検証
		if len(output.DailyStatistics) != 7 {
			t.Errorf("len(DailyStatistics) = %d, want 7", len(output.DailyStatistics))
		}
	})

	t.Run("正常系_データがない場合", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		period, _ := vo.NewStatisticsPeriod("week")

		recordRepo := &mockRecordRepository{
			getDailyCalories: func(ctx context.Context, uid vo.UserID, p vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
				return []repository.DailyCalories{}, nil
			},
		}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		output, err := uc.GetStatistics(context.Background(), userID, period)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 日数は0
		if output.TotalDays != 0 {
			t.Errorf("TotalDays = %d, want 0", output.TotalDays)
		}

		// 達成日数は0
		if output.AchievedDays != 0 {
			t.Errorf("AchievedDays = %d, want 0", output.AchievedDays)
		}

		// 超過日数は0
		if output.OverDays != 0 {
			t.Errorf("OverDays = %d, want 0", output.OverDays)
		}

		// 平均カロリーは0
		if output.AverageCalories.Value() != 0 {
			t.Errorf("AverageCalories = %d, want 0", output.AverageCalories.Value())
		}
	})

	t.Run("正常系_月間統計データを取得", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		period, _ := vo.NewStatisticsPeriod("month")

		recordRepo := &mockRecordRepository{
			getDailyCalories: func(ctx context.Context, uid vo.UserID, p vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
				// 期間がmonthであることを確認
				if p.String() != "month" {
					t.Errorf("period = %s, want month", p.String())
				}
				return []repository.DailyCalories{}, nil
			},
		}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		output, err := uc.GetStatistics(context.Background(), userID, period)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 期間の検証
		if output.Period.String() != "month" {
			t.Errorf("Period = %s, want month", output.Period.String())
		}
	})

	t.Run("正常系_平均カロリーの計算", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		period, _ := vo.NewStatisticsPeriod("week")

		// 3日分のデータ: 1000, 2000, 3000 -> 平均 2000
		now := time.Now()
		dailyCalories := []repository.DailyCalories{
			{Date: vo.ReconstructEatenAt(now.AddDate(0, 0, -2)), Calories: vo.ReconstructCalories(1000)},
			{Date: vo.ReconstructEatenAt(now.AddDate(0, 0, -1)), Calories: vo.ReconstructCalories(2000)},
			{Date: vo.ReconstructEatenAt(now), Calories: vo.ReconstructCalories(3000)},
		}

		recordRepo := &mockRecordRepository{
			getDailyCalories: func(ctx context.Context, uid vo.UserID, p vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
				return dailyCalories, nil
			},
		}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		output, err := uc.GetStatistics(context.Background(), userID, period)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// 平均カロリーの検証
		expectedAverage := 2000
		if output.AverageCalories.Value() != expectedAverage {
			t.Errorf("AverageCalories = %d, want %d", output.AverageCalories.Value(), expectedAverage)
		}
	})

	t.Run("異常系_ユーザーが存在しない", func(t *testing.T) {
		userID := vo.NewUserID()
		period, _ := vo.NewStatisticsPeriod("week")

		recordRepo := &mockRecordRepository{}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, nil
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		_, err := uc.GetStatistics(context.Background(), userID, period)

		if !errors.Is(err, domainErrors.ErrUserNotFound) {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_ユーザー取得時にエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		period, _ := vo.NewStatisticsPeriod("week")
		repoErr := errors.New("db error")

		recordRepo := &mockRecordRepository{}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return nil, repoErr
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		_, err := uc.GetStatistics(context.Background(), userID, period)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_DailyCalories取得時にエラー", func(t *testing.T) {
		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		period, _ := vo.NewStatisticsPeriod("week")
		repoErr := errors.New("db error")

		recordRepo := &mockRecordRepository{
			getDailyCalories: func(ctx context.Context, uid vo.UserID, p vo.StatisticsPeriod) ([]repository.DailyCalories, error) {
				return nil, repoErr
			},
		}
		userRepo := &mockRecordUserRepository{
			findByID: func(ctx context.Context, id vo.UserID) (*entity.User, error) {
				return user, nil
			},
		}
		recordPfcRepo := &mockRecordPfcRepository{}
		adviceCacheRepo := &mockRecordAdviceCacheRepository{}
		txManager := &mockRecordTransactionManager{}

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager)
		_, err := uc.GetStatistics(context.Background(), userID, period)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})
}
