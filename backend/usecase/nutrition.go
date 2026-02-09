package usecase

import (
	"context"
	"fmt"
	"time"

	"caltrack/config"
	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/repository"
	"caltrack/domain/vo"
	"caltrack/usecase/service"
)

const (
	// NoRecordAdviceMessage は今日の記録がない場合のアドバイスメッセージ
	NoRecordAdviceMessage = "今日の記録がまだありません。食事を記録してアドバイスを受け取りましょう！"
)

// TodayPfcOutput は今日1日のPFC摂取量と目標の出力構造体
type TodayPfcOutput struct {
	Date       time.Time // 対象日付
	CurrentPfc vo.Pfc    // 今日のPFC摂取量合計
	TargetPfc  vo.Pfc    // 目標PFC
}

// NutritionUsecase は栄養分析に関するユースケースを提供する
type NutritionUsecase struct {
	userRepo        repository.UserRepository
	recordRepo      repository.RecordRepository
	recordPfcRepo   repository.RecordPfcRepository
	adviceCacheRepo repository.AdviceCacheRepository
	pfcAnalyzer     service.PfcAnalyzer
}

// NewNutritionUsecase は NutritionUsecase のインスタンスを生成する
func NewNutritionUsecase(
	userRepo repository.UserRepository,
	recordRepo repository.RecordRepository,
	recordPfcRepo repository.RecordPfcRepository,
	adviceCacheRepo repository.AdviceCacheRepository,
	pfcAnalyzer service.PfcAnalyzer,
) *NutritionUsecase {
	return &NutritionUsecase{
		userRepo:        userRepo,
		recordRepo:      recordRepo,
		recordPfcRepo:   recordPfcRepo,
		adviceCacheRepo: adviceCacheRepo,
		pfcAnalyzer:     pfcAnalyzer,
	}
}

// GetAdvice はユーザーに対する栄養アドバイスを取得する
func (u *NutritionUsecase) GetAdvice(ctx context.Context, userID vo.UserID) (*service.NutritionAdviceOutput, error) {
	// ユーザー取得
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		logError("GetAdvice", err, "user_id", userID.String())
		return nil, err
	}
	if user == nil {
		logWarn("GetAdvice", "user not found", "user_id", userID.String())
		return nil, domainErrors.ErrUserNotFound
	}

	// 今日の日付範囲を計算
	now := time.Now()
	start := startOfDay(now)
	end := endOfDay(now)

	// 今日のRecord取得
	records, err := u.recordRepo.FindByUserIDAndDateRange(ctx, userID, start, end)
	if err != nil {
		logError("GetAdvice", err, "user_id", userID.String())
		return nil, err
	}

	// 今日の記録がない場合は固定文言を返却
	if len(records) == 0 {
		return &service.NutritionAdviceOutput{
			Advice: NoRecordAdviceMessage,
		}, nil
	}

	// キャッシュ確認
	cachedAdvice, err := u.adviceCacheRepo.FindByUserIDAndDate(ctx, userID, now)
	if err != nil {
		logError("GetAdvice", err, "user_id", userID.String())
		return nil, err
	}
	// キャッシュがあればそれを返却
	if cachedAdvice != nil {
		return &service.NutritionAdviceOutput{
			Advice: cachedAdvice.Advice(),
		}, nil
	}

	// RecordIDリスト抽出
	recordIDs := make([]vo.RecordID, len(records))
	for i, record := range records {
		recordIDs[i] = record.ID()
	}

	// 今日のRecordPfc取得
	recordPfcs := make([]*vo.Pfc, 0)
	if len(recordIDs) > 0 {
		pfcList, err := u.recordPfcRepo.FindByRecordIDs(ctx, recordIDs)
		if err != nil {
			logError("GetAdvice", err, "user_id", userID.String())
			return nil, err
		}
		for _, recordPfc := range pfcList {
			pfc := recordPfc.Pfc()
			recordPfcs = append(recordPfcs, &pfc)
		}
	}

	// 目標値計算
	targetCalories := user.CalculateTargetCalories()
	targetPfc := user.CalculateTargetPfc()

	// 現在値集計
	currentCalories := 0
	currentProtein := 0.0
	currentFat := 0.0
	currentCarbs := 0.0
	for _, record := range records {
		currentCalories += record.TotalCalories()
	}
	for _, pfc := range recordPfcs {
		currentProtein += pfc.Protein()
		currentFat += pfc.Fat()
		currentCarbs += pfc.Carbs()
	}
	currentPfc := vo.NewPfc(currentProtein, currentFat, currentCarbs)

	// 食品リスト抽出
	foodItems := make([]string, 0)
	for _, record := range records {
		foodItems = append(foodItems, record.ItemNames()...)
	}

	// 最新記録の時間帯コンテキストを取得
	latestRecord := findLatestRecord(records)
	timeContext := latestRecord.EatenAt().TimeContext()

	// PfcAnalyzer.Analyze呼び出し
	input := service.NutritionAdviceInput{
		TargetCalories:  targetCalories,
		TargetPfc:       targetPfc,
		CurrentCalories: currentCalories,
		CurrentPfc:      currentPfc,
		FoodItems:       foodItems,
		TimeContext:     timeContext,
	}

	// プロンプト構築
	prompt := buildNutritionAdvicePrompt(input)

	// Config構築
	analyzerConfig := service.PfcAnalyzerConfig{
		ModelName: config.GeminiModelName,
		Prompt:    prompt,
		Log: service.PfcAnalyzerLogConfig{
			EnableRequestLog:  true,
			EnableResponseLog: true,
			EnableTokenLog:    true,
		},
	}

	output, err := u.pfcAnalyzer.Analyze(ctx, analyzerConfig, input)
	if err != nil {
		logError("GetAdvice", err, "user_id", userID.String())
		return nil, err
	}

	// キャッシュを保存
	cache := entity.NewAdviceCache(userID, now, output.Advice)
	if err := u.adviceCacheRepo.Save(ctx, cache); err != nil {
		// キャッシュ保存失敗はログのみ（アドバイス取得は成功として扱う）
		logError("GetAdvice", err, "user_id", userID.String(), "cache_save_failed", true)
	}

	return output, nil
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

// buildNutritionAdvicePrompt は栄養アドバイスプロンプトを構築する
func buildNutritionAdvicePrompt(input service.NutritionAdviceInput) string {
	prompt := fmt.Sprintf(`あなたは栄養アドバイザーです。以下の情報に基づいて、簡潔なアドバイスを提供してください。

【時間帯情報】
%s

【目標値】
- カロリー: %d kcal
- PFC: タンパク質 %.1fg / 脂質 %.1fg / 炭水化物 %.1fg

【現在の摂取量】
- カロリー: %d kcal
- PFC: タンパク質 %.1fg / 脂質 %.1fg / 炭水化物 %.1fg

【本日食べたもの】
%s

アドバイスは以下の形式で出力してください：
- 3〜5行程度の簡潔な文章
- 目標達成度を評価
- 不足または過剰な栄養素を指摘
- 次の食事で何を意識すべきか提案`,
		input.TimeContext,
		input.TargetCalories,
		input.TargetPfc.Protein(),
		input.TargetPfc.Fat(),
		input.TargetPfc.Carbs(),
		input.CurrentCalories,
		input.CurrentPfc.Protein(),
		input.CurrentPfc.Fat(),
		input.CurrentPfc.Carbs(),
		formatFoodItems(input.FoodItems),
	)

	return prompt
}

// formatFoodItems は食品リストを整形する
func formatFoodItems(items []string) string {
	if len(items) == 0 {
		return "（まだ食事記録がありません）"
	}

	result := ""
	for i, item := range items {
		result += fmt.Sprintf("%d. %s\n", i+1, item)
	}
	return result
}

// findLatestRecord は記録リストから最新の記録を返す
func findLatestRecord(records []*entity.Record) *entity.Record {
	latest := records[0]
	for _, record := range records[1:] {
		if record.EatenAt().Time().After(latest.EatenAt().Time()) {
			latest = record
		}
	}
	return latest
}
