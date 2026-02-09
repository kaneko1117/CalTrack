package record

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"caltrack/domain/entity"
	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/common"
	"caltrack/handler/record/dto"
	"caltrack/usecase"
)

// RecordUsecaseInterface はRecordUsecaseのインターフェース
type RecordUsecaseInterface interface {
	Create(ctx context.Context, record *entity.Record) error
	GetTodayCalories(ctx context.Context, userID vo.UserID) (*usecase.TodayCaloriesOutput, error)
	GetStatistics(ctx context.Context, userID vo.UserID, period vo.StatisticsPeriod) (*usecase.StatisticsOutput, error)
}

// RecordHandler はカロリー記録関連のHTTPハンドラ
type RecordHandler struct {
	usecase RecordUsecaseInterface
}

// NewRecordHandler は RecordHandler のインスタンスを生成する
func NewRecordHandler(uc RecordUsecaseInterface) *RecordHandler {
	return &RecordHandler{usecase: uc}
}

// Create はカロリー記録を作成する
// @Summary カロリー記録作成
// @Description 食事のカロリー記録を作成する
// @Tags records
// @Accept json
// @Produce json
// @Param request body dto.CreateRecordRequest true "カロリー記録作成リクエスト"
// @Success 201 {object} dto.CreateRecordResponse "作成成功"
// @Failure 400 {object} common.ErrorResponse "リクエスト不正"
// @Failure 401 {object} common.ErrorResponse "認証失敗"
// @Failure 500 {object} common.ErrorResponse "サーバーエラー"
// @Router /records [post]
func (h *RecordHandler) Create(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userIDStr, exists := c.Get("userID")
	if !exists {
		common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "User not authenticated", nil)
		return
	}

	// リクエストボディのバインド
	var req dto.CreateRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespondError(c, http.StatusBadRequest, common.CodeInvalidRequest, "Invalid request body", nil)
		return
	}

	// 明細が空の場合はバリデーションエラー
	if len(req.Items) == 0 {
		common.RespondValidationError(c, []string{"at least one item is required"})
		return
	}

	// リクエストをEntityに変換
	record, parseErr, validationErrs := req.ToDomain(userIDStr.(string))
	if parseErr != nil {
		common.RespondError(c, http.StatusBadRequest, common.CodeValidationError, "Invalid eatenAt format", nil)
		return
	}
	if validationErrs != nil {
		details := common.ExtractErrorMessages(validationErrs)
		common.RespondValidationError(c, details)
		return
	}

	// Usecase実行
	if err := h.usecase.Create(c.Request.Context(), record); err != nil {
		common.RespondError(c, http.StatusInternalServerError, common.CodeInternalError, "Internal server error", err)
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusCreated, dto.NewCreateRecordResponse(record))
}

// GetToday は今日の摂取カロリーを取得する
// @Summary 今日の摂取カロリー取得
// @Description 認証ユーザーの今日の摂取カロリー情報を取得する
// @Tags records
// @Produce json
// @Success 200 {object} dto.TodayCaloriesResponse "取得成功"
// @Failure 401 {object} common.ErrorResponse "認証失敗"
// @Failure 404 {object} common.ErrorResponse "ユーザーが見つからない"
// @Failure 500 {object} common.ErrorResponse "サーバーエラー"
// @Router /records/today [get]
func (h *RecordHandler) GetToday(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userIDStr, exists := c.Get("userID")
	if !exists {
		common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "User not authenticated", nil)
		return
	}

	// UserID VOに変換
	userID := vo.ReconstructUserID(userIDStr.(string))

	// Usecase実行
	output, err := h.usecase.GetTodayCalories(c.Request.Context(), userID)
	if err != nil {
		// ユーザーが見つからない場合
		if errors.Is(err, domainErrors.ErrUserNotFound) {
			common.RespondError(c, http.StatusNotFound, common.CodeNotFound, "User not found", nil)
			return
		}
		common.RespondError(c, http.StatusInternalServerError, common.CodeInternalError, "Internal server error", err)
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, dto.NewTodayCaloriesResponse(output))
}

// GetStatistics は統計データを取得する
// @Summary 統計データ取得
// @Description 認証ユーザーの統計データを取得する
// @Tags records
// @Produce json
// @Param period query string false "統計期間（week または month）"
// @Success 200 {object} dto.StatisticsResponse "取得成功"
// @Failure 400 {object} common.ErrorResponse "リクエスト不正"
// @Failure 401 {object} common.ErrorResponse "認証失敗"
// @Failure 404 {object} common.ErrorResponse "ユーザーが見つからない"
// @Failure 500 {object} common.ErrorResponse "サーバーエラー"
// @Router /statistics [get]
func (h *RecordHandler) GetStatistics(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userIDStr, exists := c.Get("userID")
	if !exists {
		common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "User not authenticated", nil)
		return
	}

	// クエリパラメータのバインド
	var req dto.GetStatisticsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		common.RespondError(c, http.StatusBadRequest, common.CodeInvalidRequest, "Invalid query parameters", nil)
		return
	}

	// リクエストをVOに変換
	period, err := req.ToDomain()
	if err != nil {
		if errors.Is(err, domainErrors.ErrInvalidStatisticsPeriod) {
			common.RespondValidationError(c, []string{err.Error()})
			return
		}
		common.RespondError(c, http.StatusBadRequest, common.CodeValidationError, "Invalid period", nil)
		return
	}

	// UserID VOに変換
	userID := vo.ReconstructUserID(userIDStr.(string))

	// Usecase実行
	output, err := h.usecase.GetStatistics(c.Request.Context(), userID, period)
	if err != nil {
		if errors.Is(err, domainErrors.ErrUserNotFound) {
			common.RespondError(c, http.StatusNotFound, common.CodeNotFound, "User not found", nil)
			return
		}
		common.RespondError(c, http.StatusInternalServerError, common.CodeInternalError, "Internal server error", err)
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, dto.NewStatisticsResponse(output))
}
