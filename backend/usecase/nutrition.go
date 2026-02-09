package usecase

import (
	"context"
	"time"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/repository"
	"caltrack/domain/vo"
)

// TodayPfcOutput は今日1日のPFC摂取量と目標の出力構造体
type TodayPfcOutput struct {
	Date       time.Time // 対象日付
	CurrentPfc vo.Pfc    // 今日のPFC摂取量合計
	TargetPfc  vo.Pfc    // 目標PFC
}

// TodayPfcOutput は今日1日のPFC摂取量と目標の出力構造体
type TodayPfcOutput struct {
	Date       time.Time // 対象日付
	CurrentPfc vo.Pfc    // 今日のPFC摂取量合計
	TargetPfc  vo.Pfc    // 目標PFC
}

// NutritionUsecase は栄養分析に関するユースケースを提供する
type NutritionUsecase struct {
	userRepo      repository.UserRepository
	recordPfcRepo repository.RecordPfcRepository
}

// NewNutritionUsecase は NutritionUsecase のインスタンスを生成する
func NewNutritionUsecase(
	userRepo repository.UserRepository,
	recordPfcRepo repository.RecordPfcRepository,
) *NutritionUsecase {
	return &NutritionUsecase{
		userRepo:      userRepo,
		recordPfcRepo: recordPfcRepo,
	}
}

// GetTodayPfc は認証ユーザーの今日1日のPFC摂取量と目標を取得する
func (u *NutritionUsecase) GetTodayPfc(ctx context.Context, userID vo.UserID) (*TodayPfcOutput, error) {
	// ユーザー取得（目標PFC計算のため）
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		logError("GetTodayPfc", err, "user_id", userID.String())
		return nil, err
	}
	if user == nil {
		logWarn("GetTodayPfc", "user not found", "user_id", userID.String())
		return nil, domainErrors.ErrUserNotFound
	}

	// 今日の日付範囲を計算
	now := time.Now()
	start := startOfDay(now)
	end := endOfDay(now)

	// SQL集計でPFC合計を取得
	dailyPfc, err := u.recordPfcRepo.GetDailyPfc(ctx, userID, start, end)
	if err != nil {
		logError("GetTodayPfc", err, "user_id", userID.String())
		return nil, err
	}

	// 目標PFC計算
	targetPfc := user.CalculateTargetPfc()

	return &TodayPfcOutput{
		Date:       start,
		CurrentPfc: dailyPfc.Pfc,
		TargetPfc:  targetPfc,
	}, nil
}

// GetTodayPfc は認証ユーザーの今日1日のPFC摂取量と目標を取得する
func (u *NutritionUsecase) GetTodayPfc(ctx context.Context, userID vo.UserID) (*TodayPfcOutput, error) {
	// ユーザー取得（目標PFC計算のため）
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		logError("GetTodayPfc", err, "user_id", userID.String())
		return nil, err
	}
	if user == nil {
		logWarn("GetTodayPfc", "user not found", "user_id", userID.String())
		return nil, domainErrors.ErrUserNotFound
	}

	// 今日の日付範囲を計算
	now := time.Now()
	start := startOfDay(now)
	end := endOfDay(now)

	// SQL集計でPFC合計を取得
	dailyPfc, err := u.recordPfcRepo.GetDailyPfc(ctx, userID, start, end)
	if err != nil {
		logError("GetTodayPfc", err, "user_id", userID.String())
		return nil, err
	}

	// 目標PFC計算
	targetPfc := user.CalculateTargetPfc()

	return &TodayPfcOutput{
		Date:       start,
		CurrentPfc: dailyPfc.Pfc,
		TargetPfc:  targetPfc,
	}, nil
}
