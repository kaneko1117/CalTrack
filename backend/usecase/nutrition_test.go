package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/mock"
	"caltrack/usecase"
	"caltrack/usecase/service"

	gomock "go.uber.org/mock/gomock"
)

// setupNutritionMocks はテスト用のモックを初期化する
func setupNutritionMocks(t *testing.T) (
	*mock.MockUserRepository,
	*mock.MockRecordRepository,
	*mock.MockRecordPfcRepository,
	*mock.MockAdviceCacheRepository,
	*mock.MockPfcAnalyzer,
	*gomock.Controller,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockUserRepository(ctrl),
		mock.NewMockRecordRepository(ctrl),
		mock.NewMockRecordPfcRepository(ctrl),
		mock.NewMockAdviceCacheRepository(ctrl),
		mock.NewMockPfcAnalyzer(ctrl),
		ctrl
}

func TestNutritionUsecase_GetAdvice(t *testing.T) {
	t.Run("正常系_今日の記録がない場合は固定文言が返される", func(t *testing.T) {
		userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordRepo.EXPECT().
			FindByUserIDAndDateRange(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return([]*entity.Record{}, nil)

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
		userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		// 今日のRecordを作成
		record1, _ := entity.NewRecord(userID, time.Now())
		_ = record1.AddItem("朝食", 300)
		records := []*entity.Record{record1}

		cachedAdvice := "キャッシュされたアドバイス"
		cache := entity.NewAdviceCache(userID, time.Now(), cachedAdvice)

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordRepo.EXPECT().
			FindByUserIDAndDateRange(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return(records, nil)

		adviceCacheRepo.EXPECT().
			FindByUserIDAndDate(gomock.Any(), userID, gomock.Any()).
			Return(cache, nil)

		// analyzer.Analyzeは呼ばれないこと（EXPECTを設定しないことで検証）

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		output, err := uc.GetAdvice(context.Background(), userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.Advice != cachedAdvice {
			t.Errorf("Advice = %s, want %s", output.Advice, cachedAdvice)
		}
	})

	t.Run("正常系_キャッシュがない場合はAI呼び出し後にキャッシュ保存される", func(t *testing.T) {
		userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		// 今日のRecordを作成
		record1, _ := entity.NewRecord(userID, time.Now())
		_ = record1.AddItem("朝食", 300)

		record2, _ := entity.NewRecord(userID, time.Now())
		_ = record2.AddItem("昼食", 500)

		records := []*entity.Record{record1, record2}

		// RecordPfcを作成
		recordPfc1 := entity.NewRecordPfc(record1.ID(), 15.0, 10.0, 40.0)
		recordPfc2 := entity.NewRecordPfc(record2.ID(), 20.0, 15.0, 60.0)

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordRepo.EXPECT().
			FindByUserIDAndDateRange(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return(records, nil)

		adviceCacheRepo.EXPECT().
			FindByUserIDAndDate(gomock.Any(), userID, gomock.Any()).
			Return(nil, nil) // キャッシュなし

		recordPfcRepo.EXPECT().
			FindByRecordIDs(gomock.Any(), gomock.Any()).
			Return([]*entity.RecordPfc{recordPfc1, recordPfc2}, nil)

		analyzer.EXPECT().
			Analyze(gomock.Any(), gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, config service.PfcAnalyzerConfig, input service.NutritionAdviceInput) (*service.NutritionAdviceOutput, error) {
				// プロンプトが構築されていることを確認
				if config.Prompt == "" {
					t.Error("Prompt should not be empty")
				}
				return &service.NutritionAdviceOutput{Advice: "バランスの良い食事ができています"}, nil
			})

		adviceCacheRepo.EXPECT().
			Save(gomock.Any(), gomock.Any()).
			DoAndReturn(func(ctx context.Context, cache *entity.AdviceCache) error {
				if cache.Advice() != "バランスの良い食事ができています" {
					t.Errorf("Saved cache advice = %s, want %s", cache.Advice(), "バランスの良い食事ができています")
				}
				return nil
			})

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
	})

	t.Run("異常系_ユーザーが存在しない", func(t *testing.T) {
		userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(nil, nil)

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		_, err := uc.GetAdvice(context.Background(), userID)

		if !errors.Is(err, domainErrors.ErrUserNotFound) {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_ユーザー取得時にエラー", func(t *testing.T) {
		userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		repoErr := errors.New("db error")

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(nil, repoErr)

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		_, err := uc.GetAdvice(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_Record取得時にエラー", func(t *testing.T) {
		userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		repoErr := errors.New("db error")

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordRepo.EXPECT().
			FindByUserIDAndDateRange(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return(nil, repoErr)

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		_, err := uc.GetAdvice(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_RecordPfc取得時にエラー", func(t *testing.T) {
		userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		repoErr := errors.New("db error")

		record1, _ := entity.NewRecord(userID, time.Now())
		_ = record1.AddItem("朝食", 300)

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordRepo.EXPECT().
			FindByUserIDAndDateRange(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return([]*entity.Record{record1}, nil)

		adviceCacheRepo.EXPECT().
			FindByUserIDAndDate(gomock.Any(), userID, gomock.Any()).
			Return(nil, nil)

		recordPfcRepo.EXPECT().
			FindByRecordIDs(gomock.Any(), gomock.Any()).
			Return(nil, repoErr)

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		_, err := uc.GetAdvice(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_PfcAnalyzer実行時にエラー", func(t *testing.T) {
		userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		analyzeErr := errors.New("analyzer error")

		record1, _ := entity.NewRecord(userID, time.Now())
		_ = record1.AddItem("朝食", 300)

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordRepo.EXPECT().
			FindByUserIDAndDateRange(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return([]*entity.Record{record1}, nil)

		adviceCacheRepo.EXPECT().
			FindByUserIDAndDate(gomock.Any(), userID, gomock.Any()).
			Return(nil, nil)

		recordPfcRepo.EXPECT().
			FindByRecordIDs(gomock.Any(), gomock.Any()).
			Return([]*entity.RecordPfc{}, nil)

		analyzer.EXPECT().
			Analyze(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, analyzeErr)

		uc := usecase.NewNutritionUsecase(userRepo, recordRepo, recordPfcRepo, adviceCacheRepo, analyzer)
		_, err := uc.GetAdvice(context.Background(), userID)

		if !errors.Is(err, analyzeErr) {
			t.Errorf("got %v, want analyzeErr", err)
		}
	})
}

func TestNutritionUsecase_GetTodayPfc(t *testing.T) {
	t.Run("正常系_今日のPFC摂取量と目標を取得", func(t *testing.T) {
		userRepo, _, recordPfcRepo, _, _, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		dailyPfc := vo.NewDailyPfc(35.0, 25.0, 100.0)

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordPfcRepo.EXPECT().
			GetDailyPfc(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return(dailyPfc, nil)

		uc := usecase.NewNutritionUsecase(userRepo, nil, recordPfcRepo, nil, nil)
		output, err := uc.GetTodayPfc(context.Background(), userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.CurrentPfc.Protein() != 35.0 {
			t.Errorf("CurrentPfc.Protein() = %v, want 35.0", output.CurrentPfc.Protein())
		}
		if output.CurrentPfc.Fat() != 25.0 {
			t.Errorf("CurrentPfc.Fat() = %v, want 25.0", output.CurrentPfc.Fat())
		}
		if output.CurrentPfc.Carbs() != 100.0 {
			t.Errorf("CurrentPfc.Carbs() = %v, want 100.0", output.CurrentPfc.Carbs())
		}

		expectedTargetPfc := user.CalculateTargetPfc()
		if output.TargetPfc.Protein() != expectedTargetPfc.Protein() {
			t.Errorf("TargetPfc.Protein() = %v, want %v", output.TargetPfc.Protein(), expectedTargetPfc.Protein())
		}
		if output.TargetPfc.Fat() != expectedTargetPfc.Fat() {
			t.Errorf("TargetPfc.Fat() = %v, want %v", output.TargetPfc.Fat(), expectedTargetPfc.Fat())
		}
		if output.TargetPfc.Carbs() != expectedTargetPfc.Carbs() {
			t.Errorf("TargetPfc.Carbs() = %v, want %v", output.TargetPfc.Carbs(), expectedTargetPfc.Carbs())
		}
	})

	t.Run("正常系_記録がない場合はゼロPFCが返される", func(t *testing.T) {
		userRepo, _, recordPfcRepo, _, _, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		dailyPfc := vo.NewDailyPfc(0, 0, 0)

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordPfcRepo.EXPECT().
			GetDailyPfc(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return(dailyPfc, nil)

		uc := usecase.NewNutritionUsecase(userRepo, nil, recordPfcRepo, nil, nil)
		output, err := uc.GetTodayPfc(context.Background(), userID)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if output.CurrentPfc.Protein() != 0 {
			t.Errorf("CurrentPfc.Protein() = %v, want 0", output.CurrentPfc.Protein())
		}
		if output.CurrentPfc.Fat() != 0 {
			t.Errorf("CurrentPfc.Fat() = %v, want 0", output.CurrentPfc.Fat())
		}
		if output.CurrentPfc.Carbs() != 0 {
			t.Errorf("CurrentPfc.Carbs() = %v, want 0", output.CurrentPfc.Carbs())
		}
	})

	t.Run("異常系_ユーザーが存在しない", func(t *testing.T) {
		userRepo, _, recordPfcRepo, _, _, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(nil, nil)

		uc := usecase.NewNutritionUsecase(userRepo, nil, recordPfcRepo, nil, nil)
		_, err := uc.GetTodayPfc(context.Background(), userID)

		if !errors.Is(err, domainErrors.ErrUserNotFound) {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_ユーザー取得時にエラー", func(t *testing.T) {
		userRepo, _, recordPfcRepo, _, _, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		repoErr := errors.New("db error")

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(nil, repoErr)

		uc := usecase.NewNutritionUsecase(userRepo, nil, recordPfcRepo, nil, nil)
		_, err := uc.GetTodayPfc(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_DailyPfc取得時にエラー", func(t *testing.T) {
		userRepo, _, recordPfcRepo, _, _, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		repoErr := errors.New("db error")

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordPfcRepo.EXPECT().
			GetDailyPfc(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return(vo.DailyPfc{}, repoErr)

		uc := usecase.NewNutritionUsecase(userRepo, nil, recordPfcRepo, nil, nil)
		_, err := uc.GetTodayPfc(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})
}
