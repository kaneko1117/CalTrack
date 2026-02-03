package usecase

import (
	"context"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/repository"
	"caltrack/domain/vo"
)

// TodayCaloriesOutput は今日の摂取カロリー情報を表す出力構造体
type TodayCaloriesOutput struct {
	Date           time.Time        // 対象日付
	TotalCalories  int              // 今日の合計カロリー
	TargetCalories int              // 目標カロリー
	Difference     int              // 差分（目標 - 実績）：プラスは残り、マイナスは超過
	Records        []*entity.Record // 今日のRecord一覧
}

// RecordUsecase はカロリー記録に関するユースケースを提供する
type RecordUsecase struct {
	recordRepo repository.RecordRepository
	userRepo   repository.UserRepository
	txManager  repository.TransactionManager
}

// NewRecordUsecase は RecordUsecase のインスタンスを生成する
func NewRecordUsecase(
	recordRepo repository.RecordRepository,
	userRepo repository.UserRepository,
	txManager repository.TransactionManager,
) *RecordUsecase {
	return &RecordUsecase{
		recordRepo: recordRepo,
		userRepo:   userRepo,
		txManager:  txManager,
	}
}

// Create は新しいカロリー記録を作成する
func (u *RecordUsecase) Create(ctx context.Context, record *entity.Record) error {
	err := u.txManager.Execute(ctx, func(txCtx context.Context) error {
		if err := u.recordRepo.Save(txCtx, record); err != nil {
			logError("Create", err, "record_id", record.ID().String())
			return err
		}
		return nil
	})

	return err
}

// GetTodayCalories は認証ユーザーの今日の摂取カロリー情報を取得する
func (u *RecordUsecase) GetTodayCalories(ctx context.Context, userID vo.UserID) (*TodayCaloriesOutput, error) {
	// ユーザー取得
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		logError("GetTodayCalories", err, "user_id", userID.String())
		return nil, err
	}
	if user == nil {
		logWarn("GetTodayCalories", "user not found", "user_id", userID.String())
		return nil, domainErrors.ErrUserNotFound
	}

	// 今日の日付範囲を計算
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)

	// 今日のRecord取得
	records, err := u.recordRepo.FindByUserIDAndDateRange(ctx, userID, startOfDay, endOfDay)
	if err != nil {
		logError("GetTodayCalories", err, "user_id", userID.String())
		return nil, err
	}

	// 合計カロリー計算
	totalCalories := 0
	for _, record := range records {
		totalCalories += record.TotalCalories()
	}

	// 目標カロリー計算
	targetCalories := user.CalculateTargetCalories()

	return &TodayCaloriesOutput{
		Date:           startOfDay,
		TotalCalories:  totalCalories,
		TargetCalories: targetCalories,
		Difference:     targetCalories - totalCalories,
		Records:        records,
	}, nil
}
