package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/repository"
	"caltrack/domain/vo"
	"caltrack/mock"
	"caltrack/usecase"

	gomock "go.uber.org/mock/gomock"
)

// setupNutritionMocks はテスト用のモックを初期化する
func setupNutritionMocks(t *testing.T) (
	*mock.MockUserRepository,
	*mock.MockRecordPfcRepository,
	*gomock.Controller,
) {
	t.Helper()
	ctrl := gomock.NewController(t)
	return mock.NewMockUserRepository(ctrl),
		mock.NewMockRecordPfcRepository(ctrl),
		ctrl
}

func TestNutritionUsecase_GetTodayPfc(t *testing.T) {
	t.Run("正常系_今日のPFC摂取量と目標を取得", func(t *testing.T) {
		userRepo, recordPfcRepo, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		dailyPfc := &repository.DailyPfc{
			Date: time.Now(),
			Pfc:  vo.NewPfc(35.0, 25.0, 100.0),
		}

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordPfcRepo.EXPECT().
			GetDailyPfc(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return(dailyPfc, nil)

		uc := usecase.NewNutritionUsecase(userRepo, recordPfcRepo)
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
		userRepo, recordPfcRepo, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		dailyPfc := &repository.DailyPfc{
			Date: time.Now(),
			Pfc:  vo.NewPfc(0, 0, 0),
		}

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordPfcRepo.EXPECT().
			GetDailyPfc(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return(dailyPfc, nil)

		uc := usecase.NewNutritionUsecase(userRepo, recordPfcRepo)
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
		userRepo, recordPfcRepo, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(nil, nil)

		uc := usecase.NewNutritionUsecase(userRepo, recordPfcRepo)
		_, err := uc.GetTodayPfc(context.Background(), userID)

		if !errors.Is(err, domainErrors.ErrUserNotFound) {
			t.Errorf("got %v, want ErrUserNotFound", err)
		}
	})

	t.Run("異常系_ユーザー取得時にエラー", func(t *testing.T) {
		userRepo, recordPfcRepo, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		repoErr := errors.New("db error")

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(nil, repoErr)

		uc := usecase.NewNutritionUsecase(userRepo, recordPfcRepo)
		_, err := uc.GetTodayPfc(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})

	t.Run("異常系_DailyPfc取得時にエラー", func(t *testing.T) {
		userRepo, recordPfcRepo, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)
		repoErr := errors.New("db error")

		userRepo.EXPECT().
			FindByID(gomock.Any(), userID).
			Return(user, nil)

		recordPfcRepo.EXPECT().
			GetDailyPfc(gomock.Any(), userID, gomock.Any(), gomock.Any()).
			Return(nil, repoErr)

		uc := usecase.NewNutritionUsecase(userRepo, recordPfcRepo)
		_, err := uc.GetTodayPfc(context.Background(), userID)

		if !errors.Is(err, repoErr) {
			t.Errorf("got %v, want repoErr", err)
		}
	})
}

func TestNutritionUsecase_GetTodayPfc(t *testing.T) {
	t.Run("正常系_今日のPFC取得成功", func(t *testing.T) {
		userRepo, _, recordPfcRepo, _, _, ctrl := setupNutritionMocks(t)
		defer ctrl.Finish()

		userID := vo.NewUserID()
		user := validUserForRecord(t, userID)

		dailyPfc := vo.DailyPfc{
			Pfc: vo.NewPfc(50.0, 30.0, 150.0),
		}

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

		if output.CurrentPfc.Protein() != 50.0 {
			t.Errorf("CurrentPfc.Protein = %f, want 50.0", output.CurrentPfc.Protein())
		}
		if output.CurrentPfc.Fat() != 30.0 {
			t.Errorf("CurrentPfc.Fat = %f, want 30.0", output.CurrentPfc.Fat())
		}
		if output.CurrentPfc.Carbs() != 150.0 {
			t.Errorf("CurrentPfc.Carbs = %f, want 150.0", output.CurrentPfc.Carbs())
		}

		// 目標値の検証（BMRベース）
		targetPfc := user.CalculateTargetPfc()
		if output.TargetPfc.Protein() != targetPfc.Protein() {
			t.Errorf("TargetPfc.Protein = %f, want %f", output.TargetPfc.Protein(), targetPfc.Protein())
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

	t.Run("異常系_RecordPfc取得時にエラー", func(t *testing.T) {
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
