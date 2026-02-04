package analyze

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	domainErrors "caltrack/domain/errors"
	"caltrack/handler/analyze/dto"
	"caltrack/handler/common"
	"caltrack/usecase"
)

// AnalyzeUsecaseInterface はAnalyzeUsecaseのインターフェース（テスタビリティ向上のため）
type AnalyzeUsecaseInterface interface {
	AnalyzeImage(ctx context.Context, imageData string, mimeType string) (*usecase.AnalyzeOutput, error)
}

// AnalyzeHandler は画像解析関連のHTTPハンドラ
type AnalyzeHandler struct {
	usecase AnalyzeUsecaseInterface
}

// NewAnalyzeHandler は AnalyzeHandler のインスタンスを生成する
func NewAnalyzeHandler(uc AnalyzeUsecaseInterface) *AnalyzeHandler {
	return &AnalyzeHandler{usecase: uc}
}

// AnalyzeImage は画像から食品を解析してカロリー情報を返す
// @Summary 画像からカロリー分析
// @Description 画像に写っている食品を解析し、カロリー情報を返す
// @Tags analyze
// @Accept json
// @Produce json
// @Param request body dto.AnalyzeImageRequest true "画像解析リクエスト"
// @Success 200 {object} dto.AnalyzeImageResponse "解析成功"
// @Failure 400 {object} common.ErrorResponse "リクエスト不正"
// @Failure 401 {object} common.ErrorResponse "認証失敗"
// @Failure 500 {object} common.ErrorResponse "サーバーエラー"
// @Router /analyze-image [post]
func (h *AnalyzeHandler) AnalyzeImage(c *gin.Context) {
	// コンテキストからユーザーIDを取得（認証済みユーザーのみ使用可能）
	_, exists := c.Get("userID")
	if !exists {
		common.RespondError(c, http.StatusUnauthorized, common.CodeUnauthorized, "User not authenticated", nil)
		return
	}

	// リクエストボディのバインド
	var req dto.AnalyzeImageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.RespondError(c, http.StatusBadRequest, common.CodeInvalidRequest, "Invalid request body", nil)
		return
	}

	// Usecase実行
	output, err := h.usecase.AnalyzeImage(c.Request.Context(), req.ImageData, req.MimeType)
	if err != nil {
		// 入力バリデーションエラー
		if errors.Is(err, domainErrors.ErrImageDataRequired) {
			common.RespondValidationError(c, []string{err.Error()})
			return
		}
		if errors.Is(err, domainErrors.ErrMimeTypeRequired) {
			common.RespondValidationError(c, []string{err.Error()})
			return
		}
		// 食品が検出されなかった場合
		if errors.Is(err, domainErrors.ErrNoFoodDetected) {
			common.RespondError(c, http.StatusBadRequest, common.CodeNoFoodDetected, err.Error(), nil)
			return
		}
		// 画像解析失敗
		if errors.Is(err, domainErrors.ErrImageAnalysisFailed) {
			common.RespondError(c, http.StatusInternalServerError, common.CodeImageAnalysisFailed, "Image analysis failed", err)
			return
		}
		// その他のエラー
		common.RespondError(c, http.StatusInternalServerError, common.CodeInternalError, "Internal server error", err)
		return
	}

	// 成功レスポンス
	c.JSON(http.StatusOK, dto.NewAnalyzeImageResponse(output))
}
