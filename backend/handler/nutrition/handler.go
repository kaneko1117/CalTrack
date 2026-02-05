package nutrition

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domainErrors "caltrack/domain/errors"
	"caltrack/domain/vo"
	"caltrack/handler/common"
	"caltrack/handler/nutrition/dto"
	"caltrack/usecase"
)

// NutritionHandler は栄養分析関連のHTTPハンドラ
type NutritionHandler struct {
	usecase *usecase.NutritionUsecase
}

// NewNutritionHandler は NutritionHandler のインスタンスを生成する
func NewNutritionHandler(uc *usecase.NutritionUsecase) *NutritionHandler {
	return &NutritionHandler{usecase: uc}
}

// GetAdvice は栄養アドバイスを取得する
func (h *NutritionHandler) GetAdvice(c *gin.Context) {
	// コンテキストからユーザーIDを取得
	userIDStr, exists := c.Get("userID")
	if !exists {
		common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "User not authenticated", nil)
		return
	}

	// UserID VOに変換
	userID := vo.ReconstructUserID(userIDStr.(string))

	// Usecase実行
	output, err := h.usecase.GetAdvice(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, domainErrors.ErrUserNotFound) {
			common.RespondError(c, http.StatusNotFound, common.CodeNotFound, "User not found", nil)
			return
		}
		common.RespondError(c, http.StatusInternalServerError, common.CodeInternalError, "Internal server error", err)
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, dto.NewAdviceResponse(output))
}
