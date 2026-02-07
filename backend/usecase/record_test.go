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
	"caltrack/mock"
	"caltrack/usecase"
	"caltrack/usecase/service"

	gomock "go.uber.org/mock/gomock"
)

// setupRecordMocks はテスト用のモックを初期化する
func setupRecordMocks(t *testing.T) (
	*mock.MockRecordRepository,
	*mock.MockRecordPfcRepository,
	*mock.MockUserRepository,
	*mock.MockAdviceCacheRepository,
	*mock.MockTransactionManager,
	*mock.MockPfcEstimator,
	*gomock.Controller,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockRecordRepository(ctrl),
		mock.NewMockRecordPfcRepository(ctrl),
		mock.NewMockUserRepository(ctrl),
		mock.NewMockAdviceCacheRepository(ctrl),
		mock.NewMockTransactionManager(ctrl),
		mock.NewMockPfcEstimator(ctrl),
		ctrl
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

func TestRecordUsecase_Create(t *testing.T) {
	t.Run("正常系_記録が保存されキャッシュが無効化される", func(t *testing.T) {
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		record := validRecord(t)
		var savedRecord *entity.Record
		var savedRecordPfc *entity.RecordPfc
		cacheDeleted := false

		setupTxManagerExecute(txManager)
		recordRepo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, r *entity.Record) error {
				savedRecord = r
				return nil
			})
		recordPfcRepo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, rp *entity.RecordPfc) error {
				savedRecordPfc = rp
				return nil
			})
		pfcEstimator.EXPECT().
			Estimate(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(&service.PfcEstimateOutput{
				Protein: 20.0,
				Fat:     10.0,
				Carbs:   30.0,
			}, nil)
		adviceCacheRepo.EXPECT().
			DeleteByUserIDAndDate(gomock.Any(), gomock.Eq(record.UserID()), gomock.Any()).
			DoAndReturn(func(ctx context.Context, userID vo.UserID, date time.Time) error {
				cacheDeleted = true
				return nil
			})

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
		err := uc.Create(context.Background(), record)

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
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		record := validRecord(t)
		saveErr := errors.New("save error")

		setupTxManagerExecute(txManager)
		recordRepo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			Return(saveErr)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
		err := uc.Create(context.Background(), record)

		if !errors.Is(err, saveErr) {
			t.Errorf("got %v, want saveErr", err)
		}
	})
}

func TestRecordUsecase_GetTodayCalories(t *testing.T) {
	t.Run("正常系_今日のカロリー情報を取得", func(t *testing.T) {
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		// 今日のRecordを作成
		record1, _ := entity.NewRecord(userID, time.Now())
		_ = record1.AddItem("朝食", 300)

		record2, _ := entity.NewRecord(userID, time.Now())
		_ = record2.AddItem("昼食", 500)
		_ = record2.AddItem("おやつ", 100)

		records := []*entity.Record{record1, record2}

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(user, nil)
		recordRepo.EXPECT().
			FindByUserIDAndDateRange(gomock.Any(), gomock.Eq(userID), gomock.Any(), gomock.Any()).
			Return(records, nil)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
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
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(user, nil)
		recordRepo.EXPECT().
			FindByUserIDAndDateRange(gomock.Any(), gomock.Eq(userID), gomock.Any(), gomock.Any()).
			Return([]*entity.Record{}, nil)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
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
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(nil, nil)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
		_, err := uc.GetTodayCalories(context.Background(), userID)

		if !errors.Is(err, domainErrors.ErrUserNotFound) {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_ユーザー取得時にエラー", func(t *testing.T) {
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		repoErr := errors.New("db error")

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(nil, repoErr)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
		_, err := uc.GetTodayCalories(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_Record取得時にエラー", func(t *testing.T) {
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		repoErr := errors.New("db error")

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(user, nil)
		recordRepo.EXPECT().
			FindByUserIDAndDateRange(gomock.Any(), gomock.Eq(userID), gomock.Any(), gomock.Any()).
			Return(nil, repoErr)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
		_, err := uc.GetTodayCalories(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})
}

func TestRecordUsecase_GetStatistics(t *testing.T) {
	t.Run("正常系_週間統計データを取得", func(t *testing.T) {
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

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

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(user, nil)
		recordRepo.EXPECT().
			GetDailyCalories(gomock.Any(), gomock.Eq(userID), gomock.Eq(period)).
			Return(dailyCalories, nil)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
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
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		period, _ := vo.NewStatisticsPeriod("week")

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(user, nil)
		recordRepo.EXPECT().
			GetDailyCalories(gomock.Any(), gomock.Eq(userID), gomock.Eq(period)).
			Return([]repository.DailyCalories{}, nil)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
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
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		period, _ := vo.NewStatisticsPeriod("month")

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(user, nil)
		recordRepo.EXPECT().
			GetDailyCalories(gomock.Any(), gomock.Eq(userID), gomock.Eq(period)).
			Return([]repository.DailyCalories{}, nil)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
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
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

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

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(user, nil)
		recordRepo.EXPECT().
			GetDailyCalories(gomock.Any(), gomock.Eq(userID), gomock.Eq(period)).
			Return(dailyCalories, nil)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
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
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		period, _ := vo.NewStatisticsPeriod("week")

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(nil, nil)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
		_, err := uc.GetStatistics(context.Background(), userID, period)

		if !errors.Is(err, domainErrors.ErrUserNotFound) {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_ユーザー取得時にエラー", func(t *testing.T) {
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		period, _ := vo.NewStatisticsPeriod("week")
		repoErr := errors.New("db error")

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(nil, repoErr)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
		_, err := uc.GetStatistics(context.Background(), userID, period)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_DailyCalories取得時にエラー", func(t *testing.T) {
		recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator, ctrl := setupRecordMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		period, _ := vo.NewStatisticsPeriod("week")
		repoErr := errors.New("db error")

		userRepo.EXPECT().
			FindByID(gomock.Any(), gomock.Eq(userID)).
			Return(user, nil)
		recordRepo.EXPECT().
			GetDailyCalories(gomock.Any(), gomock.Eq(userID), gomock.Eq(period)).
			Return(nil, repoErr)

		uc := usecase.NewRecordUsecase(recordRepo, recordPfcRepo, userRepo, adviceCacheRepo, txManager, pfcEstimator)
		_, err := uc.GetStatistics(context.Background(), userID, period)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})
}
