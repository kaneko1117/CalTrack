package nutrition

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/common"
	"caltrack/handler/nutrition/dto"
	"caltrack/usecase"
)

// NutritionUsecaseInterface はNutritionUsecaseのインターフェース
type NutritionUsecaseInterface interface {
	GetTodayPfc(ctx context.Context, userID vo.UserID) (*usecase.TodayPfcOutput, error)
}

// NutritionHandler は栄養分析関連のHTTPハンドラ
type NutritionHandler struct {
	usecase NutritionUsecaseInterface
}

// NewNutritionHandler は NutritionHandler のインスタンスを生成する
func NewNutritionHandler(uc NutritionUsecaseInterface) *NutritionHandler {
	return &NutritionHandler{usecase: uc}
}

// GetTodayPfc は今日のPFC摂取量と目標PFCを取得する
func (h *NutritionHandler) GetTodayPfc(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userIDStr, exists := c.Get("userID")
	if !exists {
		common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "User not authenticated", nil)
		return
	}

	// UserID VOに変換
	userID := vo.ReconstructUserID(userIDStr.(string))

	// Usecase実行
	output, err := h.usecase.GetTodayPfc(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrUserNotFound) {
			common.RespondError(c, http.StatusNotFound, common.CodeNotFound, "User not found", nil)
			return
		}
		common.RespondError(c, http.StatusInternalServerError, common.CodeInternalError, "Internal server error", err)
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, dto.NewTodayPfcResponse(output))
}
