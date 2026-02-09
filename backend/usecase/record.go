package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/repository"
	"caltrack/domain/vo"
	"caltrack/usecase/service"
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
	recordRepo      repository.RecordRepository
	recordPfcRepo   repository.RecordPfcRepository
	userRepo        repository.UserRepository
	adviceCacheRepo repository.AdviceCacheRepository
	txManager       repository.TransactionManager
	pfcEstimator    service.PfcEstimator
	aiConfig        AIConfig
}

// NewRecordUsecase は RecordUsecase のインスタンスを生成する
func NewRecordUsecase(
	recordRepo repository.RecordRepository,
	recordPfcRepo repository.RecordPfcRepository,
	userRepo repository.UserRepository,
	adviceCacheRepo repository.AdviceCacheRepository,
	txManager repository.TransactionManager,
	pfcEstimator service.PfcEstimator,
	aiConfig AIConfig,
) *RecordUsecase {
	return &RecordUsecase{
		recordRepo:      recordRepo,
		recordPfcRepo:   recordPfcRepo,
		userRepo:        userRepo,
		adviceCacheRepo: adviceCacheRepo,
		txManager:       txManager,
		pfcEstimator:    pfcEstimator,
		aiConfig:        aiConfig,
	}
}

// Create は新しいカロリー記録を作成する
func (u *RecordUsecase) Create(ctx context.Context, record *entity.Record) error {
	err := u.txManager.Execute(ctx, func(txCtx context.Context) error {
		// Record保存
		if err := u.recordRepo.Save(txCtx, record); err != nil {
			logError("Create", err, "record_id", record.ID().String())
			return err
		}

		// AI-PFC推定実行
		recordPfc, err := u.estimatePfc(txCtx, record)
		if err != nil {
			// PFC推定失敗してもRecordは保存済みなのでログのみ
			logError("Create", err, "record_id", record.ID().String(), "pfc_estimation_failed", true)
		}

		// RecordPfc保存（推定成功時のみ）
		if recordPfc != nil {
			if err := u.recordPfcRepo.Save(txCtx, recordPfc); err != nil {
				logError("Create", err, "record_pfc_id", recordPfc.ID().String())
				return err
			}
		}

		// キャッシュ無効化（記録日のキャッシュを削除）
		recordDate := record.EatenAt().Time()
		if err := u.adviceCacheRepo.DeleteByUserIDAndDate(txCtx, record.UserID(), recordDate); err != nil {
			// キャッシュ削除失敗はログのみ（記録作成は成功として扱う）
			logError("Create", err, "user_id", record.UserID().String(), "cache_delete_failed", true)
		}

		return nil
	})

	return err
}

// estimatePfc は食品名からPFC値を推定してRecordPfcを作成する
func (u *RecordUsecase) estimatePfc(ctx context.Context, record *entity.Record) (*entity.RecordPfc, error) {
	// 食品名リストを抽出
	foodNames := record.ItemNames()

	// PFC推定プロンプト構築
	prompt := buildPfcEstimatePrompt(foodNames)

	// PFC推定実行
	estimatorConfig := service.PfcEstimatorConfig{
		ModelName: u.aiConfig.GeminiModelName(),
		Prompt:    prompt,
		Log: service.PfcEstimatorLogConfig{
			EnableRequestLog:  true,
			EnableResponseLog: true,
			EnableTokenLog:    true,
		},
	}

	input := service.PfcEstimateInput{
		FoodItems: foodNames,
	}

	output, err := u.pfcEstimator.Estimate(ctx, estimatorConfig, input)
	if err != nil {
		return nil, err
	}

	// RecordPfcを生成（エラーなし）
	recordPfc := entity.NewRecordPfc(
		record.ID(),
		output.Protein,
		output.Fat,
		output.Carbs,
	)

	return recordPfc, nil
}

// buildPfcEstimatePrompt はPFC推定用のプロンプトを構築する
func buildPfcEstimatePrompt(foodNames []string) string {
	foodList := strings.Join(foodNames, "\n- ")
	return fmt.Sprintf(`以下の食品リストから、合計のPFC（タンパク質・脂質・炭水化物）をグラム単位で推定してください。

食品リスト:
- %s

回答は以下のJSON形式で返してください（他の説明は不要です）:
{
  "protein": 数値,
  "fat": 数値,
  "carbs": 数値
}

例:
{
  "protein": 25.5,
  "fat": 12.3,
  "carbs": 45.0
}`, foodList)
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
	start := startOfDay(now)
	end := endOfDay(now)

	// 今日のRecord取得
	records, err := u.recordRepo.FindByUserIDAndDateRange(ctx, userID, start, end)
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
		Date:           start,
		TotalCalories:  totalCalories,
		TargetCalories: targetCalories,
		Difference:     targetCalories - totalCalories,
		Records:        records,
	}, nil
}

// DailyStatistics は日別統計データ（グラフ表示用）
type DailyStatistics struct {
	Date           vo.EatenAt  // 対象日付
	TotalCalories  vo.Calories // その日の合計カロリー
	TargetCalories vo.Calories // 目標カロリー
	IsAchieved     bool        // 達成フラグ（80%〜100%）
	IsOver         bool        // 超過フラグ（100%超）
}

// StatisticsOutput は統計データ出力
type StatisticsOutput struct {
	Period          vo.StatisticsPeriod // 統計期間（week/month）
	TargetCalories  vo.Calories         // 1日の目標カロリー
	AverageCalories vo.Calories         // 期間内の平均カロリー
	TotalDays       int                 // 期間の日数
	AchievedDays    int                 // 達成日数（80%〜100%）
	OverDays        int                 // 超過日数（100%超）
	DailyStatistics []DailyStatistics   // 日別統計データ（グラフ用）
}

// GetStatistics は認証ユーザーの統計データを取得する
func (u *RecordUsecase) GetStatistics(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) (*StatisticsOutput, error) {
	// ユーザー取得（目標カロリー計算のため）
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		logError("GetStatistics", err, "user_id", userID.String())
		return nil, err
	}
	if user == nil {
		logWarn("GetStatistics", "user not found", "user_id", userID.String())
		return nil, domainErrors.ErrUserNotFound
	}

	// 目標カロリー取得
	targetCalories := vo.ReconstructCalories(user.CalculateTargetCalories())

	// 日別カロリーデータ取得
	dailyCaloriesList, err := u.recordRepo.GetDailyCalories(ctx, userID, period)
	if err != nil {
		logError("GetStatistics", err, "user_id", userID.String())
		return nil, err
	}

	// 集計用変数の初期化
	totalDays := len(dailyCaloriesList)
	achievedDays := 0
	overDays := 0
	totalCaloriesSum := vo.ZeroCalories()
	dailyStatistics := make([]DailyStatistics, 0, totalDays)

	// 日別データをループして集計
	for _, daily := range dailyCaloriesList {
		// 達成・超過判定（VOのメソッドで判定）
		isAchieved := daily.Calories.IsAchieved(targetCalories)
		isOver := daily.Calories.IsOver(targetCalories)

		if isAchieved {
			achievedDays++
		}
		if isOver {
			overDays++
		}

		totalCaloriesSum = totalCaloriesSum.Add(daily.Calories)

		dailyStatistics = append(dailyStatistics, DailyStatistics{
			Date:           daily.Date,
			TotalCalories:  daily.Calories,
			TargetCalories: targetCalories,
			IsAchieved:     isAchieved,
			IsOver:         isOver,
		})
	}

	// 平均カロリー計算（0除算防止）
	averageCalories := vo.ZeroCalories()
	if totalDays > 0 {
		averageCalories = vo.ReconstructCalories(totalCaloriesSum.Value() / totalDays)
	}

	return &StatisticsOutput{
		Period:          period,
		TargetCalories:  targetCalories,
		AverageCalories: averageCalories,
		TotalDays:       totalDays,
		AchievedDays:    achievedDays,
		OverDays:        overDays,
		DailyStatistics: dailyStatistics,
	}, nil
}
